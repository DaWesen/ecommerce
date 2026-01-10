package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce/order-service/internal/client"
	"ecommerce/order-service/internal/dao/dao"
	"ecommerce/order-service/internal/service"
	"ecommerce/order-service/kitex_gen/api/orderservice"
	"ecommerce/order-service/pkg/config"
	"ecommerce/order-service/pkg/database"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitexServer "github.com/cloudwego/kitex/server"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("💥 加载配置失败: %v 💥", err)
	}

	log.Printf("✅ 配置加载成功")
	log.Printf("📊 数据库配置: %s@%s:%d/%s",
		cfg.Database.MySQL.User,
		cfg.Database.MySQL.Host,
		cfg.Database.MySQL.Port,
		cfg.Database.MySQL.DBName,
	)

	// 初始化数据库
	db, dbType, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("💥 数据库连接失败: %v 💥", err)
	}
	log.Printf("✅ 数据库连接成功，类型: %s", dbType)

	// 初始化服务依赖
	orderService, err := initOrderService(cfg, db)
	if err != nil {
		log.Fatalf("💥 初始化订单服务失败: %v 💥", err)
	}

	// 创建信号通道用于关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动 HTTP 服务器
	var httpServer *server.Hertz
	if cfg.Hertz.Port > 0 {
		httpServer = startHTTPServer(cfg.Hertz.Port, orderService)
	}

	// 启动 Kitex RPC 服务器
	kitexPort := cfg.Kitex.Port
	if kitexPort <= 0 {
		kitexPort = 50053
	}
	kitexServer := startKitexServer(kitexPort, orderService)

	log.Printf("🚀 服务启动成功!")
	log.Printf("🌐 HTTP API: http://localhost:%d", cfg.Hertz.Port)
	log.Printf("🔌 RPC 服务: localhost:%d", kitexPort)
	log.Printf("⚙️  模式: %s", cfg.Hertz.Mode)

	// 等待关闭信号
	<-quit
	log.Println("⏳ 收到关闭信号，开始关闭...")

	//关闭
	gracefulShutdown(httpServer, kitexServer, db)
}

// initOrderService 初始化订单服务
func initOrderService(cfg *config.Config, db *gorm.DB) (*service.OrderService, error) {
	//初始化 DAO 工厂
	daoFactory := dao.NewDaoFactory(db)

	//初始化用户服务客户端
	userServiceAddr := "127.0.0.1:50052"
	log.Printf("尝试连接用户服务: %s", userServiceAddr)

	userClient, err := client.NewUserClient(userServiceAddr)
	if err != nil {
		log.Printf("⚠️  创建用户服务客户端失败: %v", err)
		return nil, fmt.Errorf("用户服务客户端初始化失败: %v", err)
	}

	// 测试用户服务连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试用户服务是否可用
	testUserID := int64(1)
	_, err = userClient.GetUserInfo(ctx, testUserID)
	if err != nil {
		log.Printf("❌ 用户服务连接测试失败: %v", err)
		log.Println("请检查用户服务:")
		log.Println("1. 用户服务是否正在运行?")
		log.Println("2. 用户服务是否在端口 50052?")
		log.Println("3. 用户服务的服务名是否正确?")

		// 显示当前运行的服务端口
		log.Println("当前运行的服务:")
		log.Println("  从日志看，用户服务在 50052")
		log.Println("  商品服务在 50051")

		return nil, fmt.Errorf("用户服务不可用: %v", err)
	}

	log.Println("✅ 用户服务连接测试成功")
	//初始化商品
	productServiceAddr := "127.0.0.1:50051"
	log.Printf("尝试连接商品服务: %s", productServiceAddr)

	productClient, err := client.NewProductClient(productServiceAddr)
	if err != nil {
		log.Printf("⚠️  创建商品服务客户端失败: %v", err)
		return nil, fmt.Errorf("商品服务客户端初始化失败: %v", err)
	}

	// 测试商品服务连接
	_, err = productClient.GetProductInfo(ctx, 1)
	if err != nil {
		log.Printf("❌ 商品服务连接测试失败: %v", err)
		log.Println("请检查商品服务:")
		log.Println("1. 商品服务是否正在运行?")
		log.Println("2. 商品服务是否在端口 50051?")
		log.Println("3. 商品服务的服务名是否正确?")
		return nil, fmt.Errorf("商品服务不可用: %v", err)
	}

	log.Println("✅ 商品服务连接测试成功")

	//创建订单服务
	orderService := service.NewOrderService(db, daoFactory, userClient, productClient)

	log.Println("✅ 订单服务初始化成功")
	return orderService, nil
}

// startHTTPServer 启动 HTTP 服务器
func startHTTPServer(port int, orderService *service.OrderService) *server.Hertz {
	h := server.New(
		server.WithHostPorts(fmt.Sprintf(":%d", port)),
		server.WithMaxRequestBodySize(10*1024*1024),
	)

	// 注册 HTTP 路由
	registerHTTPRoutes(h, orderService)

	// 启动 HTTP 服务器
	go func() {
		log.Printf("🌐 HTTP 服务器启动在端口 %d", port)
		if err := h.Run(); err != nil {
			log.Fatalf("💥 HTTP 服务器启动失败: %v", err)
		}
	}()

	return h
}

// startKitexServer 启动 Kitex RPC 服务器
func startKitexServer(port int, orderService *service.OrderService) kitexServer.Server {
	// 创建 Kitex handler
	kitexHandler := NewOrderServiceImpl(orderService)

	// 创建服务器地址
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("💥 解析地址失败: %v", err)
	}

	// 创建 Kitex 服务器
	svr := orderservice.NewServer(
		kitexHandler,
		kitexServer.WithServiceAddr(addr),
		kitexServer.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "order.service",
		}),
	)

	// 启动 Kitex 服务器
	go func() {
		log.Printf("🔌 Kitex RPC 服务器启动在端口 %d", port)
		if err := svr.Run(); err != nil {
			log.Fatalf("💥 Kitex 服务器启动失败: %v", err)
		}
	}()

	return svr
}

// registerHTTPRoutes 注册 HTTP 路由
func registerHTTPRoutes(h *server.Hertz, orderService *service.OrderService) {
	// 健康检查
	h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"status":  "ok",
			"time":    time.Now().Unix(),
			"service": "order-service",
		})
	})

	// API 文档
	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"service":     "Order Service",
			"description": "订单微服务",
			"apis": []utils.H{
				{"method": "GET", "path": "/health", "desc": "健康检查"},
				{"method": "GET", "path": "/stats", "desc": "订单统计"},
				{"method": "GET", "path": "/orders", "desc": "查询订单列表"},
				{"method": "GET", "path": "/orders/:orderNo", "desc": "获取订单详情"},
				{"method": "POST", "path": "/orders", "desc": "创建订单"},
				{"method": "POST", "path": "/orders/:orderNo/pay", "desc": "支付订单"},
				{"method": "POST", "path": "/orders/:orderNo/cancel", "desc": "取消订单"},
				{"method": "POST", "path": "/orders/:orderNo/ship", "desc": "发货"},
				{"method": "POST", "path": "/orders/:orderNo/receive", "desc": "确认收货"},
				{"method": "POST", "path": "/orders/:orderNo/refund", "desc": "申请退款"},
			},
		})
	})

	// 订单统计
	h.GET("/stats", func(c context.Context, ctx *app.RequestContext) {
		// 这里可以调用 orderService 的统计方法
		ctx.JSON(consts.StatusOK, utils.H{
			"message": "订单统计接口",
		})
	})

	// 订单相关 API（示例，实际需要实现具体的处理函数）
	h.GET("/orders", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"message": "查询订单列表",
		})
	})

	h.GET("/orders/:orderNo", func(c context.Context, ctx *app.RequestContext) {
		orderNo := ctx.Param("orderNo")
		ctx.JSON(consts.StatusOK, utils.H{
			"message":  "获取订单详情",
			"order_no": orderNo,
		})
	})

	h.POST("/orders", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"message": "创建订单",
		})
	})
}

// 关闭
func gracefulShutdown(httpServer *server.Hertz, kitexServer kitexServer.Server, db *gorm.DB) {
	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("⏳ 开始关闭服务...")

	// 关闭 HTTP 服务器
	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("⚠️ 关闭 HTTP 服务器失败: %v", err)
		} else {
			log.Println("✅ HTTP 服务器已关闭")
		}
	}

	// 关闭 Kitex 服务器
	if kitexServer != nil {
		if err := kitexServer.Stop(); err != nil {
			log.Printf("⚠️ 关闭 Kitex 服务器失败: %v", err)
		} else {
			log.Println("✅ Kitex 服务器已关闭")
		}
	}

	// 关闭数据库连接
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("⚠️ 关闭数据库连接失败: %v", err)
			} else {
				log.Println("✅ 数据库连接已关闭")
			}
		}
	}

	log.Println("🎉 服务关闭完成")
	os.Exit(0)
}

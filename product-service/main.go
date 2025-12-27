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

	"ecommerce/product-service/internal/handler"
	"ecommerce/product-service/internal/repository"
	"ecommerce/product-service/internal/service"
	api "ecommerce/product-service/kitex_gen/api/productservice"
	"ecommerce/product-service/pkg/config"
	"ecommerce/product-service/pkg/database"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	kitexServer "github.com/cloudwego/kitex/server"
	"gorm.io/gorm"
)

func main() {
	//加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("💥加载配置失败: %v💥", err)
	}

	log.Printf("💖配置加载成功: HTTP端口=%d💖", cfg.Hertz.Port)
	log.Printf("😎数据库配置: %s@%s:%d/%s😎",
		cfg.Database.MySQL.User,
		cfg.Database.MySQL.Host,
		cfg.Database.MySQL.Port,
		cfg.Database.MySQL.DBName,
	)

	//初始化数据库
	db, _, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("💥数据库连接失败: %v💥", err)
	}

	//初始化依赖
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)

	//创建信号通道用于关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//启动Hertz
	httpServer := startHTTPServer(cfg.Hertz.Port, productService)

	//启动Kitex
	kitexPort := 50051
	if cfg.Kitex.Port > 0 {
		kitexPort = cfg.Kitex.Port
	}
	kitexServer := startKitexServer(kitexPort, productService)

	log.Printf("💖服务启动成功!💖")
	log.Printf("💖HTTP API: http://localhost:%d💖", cfg.Hertz.Port)
	log.Printf("💖RPC服务: localhost:%d💖", kitexPort)
	log.Printf("💖模式: %s💖", cfg.Hertz.Mode)

	//等待关闭信号
	<-quit
	log.Println("💖收到关闭信号，开始关闭...💖")

	//关闭
	gracefulShutdown(httpServer, kitexServer, db)
}

// 启动Hertz
func startHTTPServer(port int, productService service.ProductService) *server.Hertz {
	h := server.New(
		server.WithHostPorts(fmt.Sprintf(":%d", port)),
		server.WithMaxRequestBodySize(10*1024*1024),
	)
	// 创建HTTPhandler
	httpHandler := handler.NewProductHTTPHandler(productService)

	// 注册路由
	registerHTTPRoutes(h, httpHandler)

	//启动
	go func() {
		log.Printf("💖HTTP服务器启动在端口 %d💖", port)
		h.Spin()
	}()
	return h
}

// 启动Kitex
func startKitexServer(port int, productService service.ProductService) kitexServer.Server {
	//创建 Kitex handler
	kitexHandler := &ProductServiceImpl{
		productService: productService,
	}

	//创建服务器地址
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("😭解析地址失败: %v😭", err)
	}

	//创建Kitex
	svr := api.NewServer(
		kitexHandler,
		kitexServer.WithServiceAddr(addr),
		kitexServer.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "product.service",
		}),
	)

	//启动
	go func() {
		log.Printf("💖Kitex RPC 服务器启动在端口 %d💖", port)
		if err := svr.Run(); err != nil {
			log.Fatalf("😭Kitex 服务器启动失败: %v😭", err)
		}
	}()

	return svr
}

// 注册 HTTP 路由
func registerHTTPRoutes(h *server.Hertz, httpHandler *handler.ProductHTTPHandler) {
	//健康检查
	h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"status":  "ok",
			"time":    time.Now().Unix(),
			"service": "product-service",
		})
	})

	//API 文档
	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"service":     "Product Service",
			"description": "产品微服务",
			"apis": []utils.H{
				{"method": "POST", "path": "/api/v1/products", "desc": "创建产品"},
				{"method": "GET", "path": "/api/v1/products/:id", "desc": "获取产品详情"},
				{"method": "PUT", "path": "/api/v1/products/:id", "desc": "更新产品"},
				{"method": "DELETE", "path": "/api/v1/products/:id", "desc": "删除产品"},
				{"method": "GET", "path": "/api/v1/products", "desc": "用户搜索产品"},
				{"method": "GET", "path": "/api/v1/admin/products", "desc": "管理员搜索产品"},
				{"method": "POST", "path": "/api/v1/products/:id/online", "desc": "上架产品"},
				{"method": "POST", "path": "/api/v1/products/:id/offline", "desc": "下架产品"},
			},
		})
	})

	//产品相关 API
	h.POST("/api/v1/products", httpHandler.CreateProduct)
	h.GET("/api/v1/products/:id", httpHandler.GetProduct)
	h.PUT("/api/v1/products/:id", httpHandler.UpdateProduct)
	h.DELETE("/api/v1/products/:id", httpHandler.DeleteProduct)
	h.GET("/api/v1/products", httpHandler.SearchProducts)
	h.GET("/api/v1/admin/products", httpHandler.AdminSearchProducts)
	h.POST("/api/v1/products/:id/online", httpHandler.OnlineProduct)
	h.POST("/api/v1/products/:id/offline", httpHandler.OfflineProduct)
}

// 关闭
func gracefulShutdown(httpServer *server.Hertz, kitexServer kitexServer.Server, db *gorm.DB) {
	//创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("💖开始关闭...💖")

	//关闭 HTTP 服务器
	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("😭关闭 HTTP 服务器失败: %v😭", err)
		} else {
			log.Println("💖HTTP 服务器已关闭💖")
		}
	}

	//关闭 Kitex 服务器
	if kitexServer != nil {
		if err := kitexServer.Stop(); err != nil {
			log.Printf("😭关闭 Kitex 服务器失败: %v😭", err)
		} else {
			log.Println("💖Kitex 服务器已关闭💖")
		}
	}

	//关闭数据库连接
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("😭关闭数据库连接失败: %v😭", err)
			} else {
				log.Println("💖数据库连接已关闭💖")
			}
		}
	}
	log.Println("💖服务关闭完成💖")
	os.Exit(0)
}

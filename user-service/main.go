package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"ecommerce/user-service/internal/handler"
	"ecommerce/user-service/internal/repository"
	"ecommerce/user-service/internal/service"
	"ecommerce/user-service/kitex_gen/api/userservice"
	"ecommerce/user-service/pkg/config"
	"ecommerce/user-service/pkg/database"
	"ecommerce/user-service/pkg/jwt"
	"ecommerce/user-service/pkg/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	hertzServer "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
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

	//初始化JWT管理器
	jwtCfg := jwt.Config{
		SecretKey:     cfg.JWT.Secret,
		Issuer:        "ecommerce-user-service",
		AccessExpire:  time.Duration(cfg.JWT.ExpireHours) * time.Hour,
		RefreshExpire: 7 * 24 * time.Hour,
		Algorithm:     "HS256",
	}

	jwtManager := jwt.NewJWTManager(jwtCfg)

	//初始化依赖
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)

	//初始化认证中间件
	authMiddleware := middleware.NewAuthMiddleware(
		jwtManager,
		[]string{
			"/health",
			"/",
			"/api/v1/auth/register",
			"/api/v1/auth/login",
		},
	)

	//创建信号通道用于关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	httpServer := startHTTPServer(cfg.Hertz.Port, userService, authMiddleware)
	kitexPort := 50052
	if cfg.Kitex.Port > 0 {
		kitexPort = cfg.Kitex.Port
	}

	// 创建Kitex处理器
	kitexHandler, err := NewUserServiceImpl()
	if err != nil {
		log.Fatalf("💥创建Kitex处理器失败: %v💥", err)
	}

	kitexServer := startKitexServer(kitexPort, kitexHandler) // 只传两个参
	log.Printf("💖服务启动成功!💖")
	log.Printf("💖HTTP API: http://localhost:%d💖", cfg.Hertz.Port)
	log.Printf("💖RPC服务: localhost:%d💖", kitexPort)
	log.Printf("💖模式: %s💖", cfg.Hertz.Mode)

	//等待关闭信号
	<-quit
	log.Println("💖收到关闭信号，开始关闭...💖")

	//执行关闭
	gracefulShutdown(httpServer, kitexServer, db)
}

// 启动HTTP服务器
func startHTTPServer(port int, userService service.UserService, authMiddleware *middleware.AuthMiddleware) *hertzServer.Hertz {
	h := hertzServer.New(
		hertzServer.WithHostPorts(fmt.Sprintf(":%d", port)),
		hertzServer.WithMaxRequestBodySize(10*1024*1024),
	)

	//创建HTTP处理器
	httpHandler := handler.NewUserHTTPHandler(userService)

	//注册路由
	registerHTTPRoutes(h, httpHandler, authMiddleware)

	//关闭
	go func() {
		log.Printf("💖HTTP服务器启动在端口 %d💖", port)
		err := h.Run()
		if err != nil {
			//检查是否是因为Shutdown引起的错误
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("💖HTTP服务器已关闭💖")
				return
			}
			log.Fatalf("💥HTTP服务器启动失败: %v💥", err)
		}
	}()

	return h
}

func startKitexServer(port int, kitexHandler *UserServiceImpl) kitexServer.Server {
	//解析服务器地址
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("😭解析地址失败: %v😭", err)
	}

	//创建Kitex服务器 - 使用 TTHeader 协议
	svr := userservice.NewServer(
		kitexHandler,
		kitexServer.WithServiceAddr(addr),
		kitexServer.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "user.service",
		}),
		kitexServer.WithMetaHandler(transmeta.ServerTTHeaderHandler), //添加TTHeader处理器
		kitexServer.WithLimit(&limit.Option{
			MaxConnections: 1000,
			MaxQPS:         500,
		}),
	)

	//启动Kitex服务器
	go func() {
		log.Printf("💖Kitex RPC服务器启动在端口 %d💖", port)
		if err := svr.Run(); err != nil {
			klog.Fatalf("😭Kitex服务器启动失败: %v😭", err)
		}
	}()

	return svr
}

// 注册HTTP路由
func registerHTTPRoutes(h *hertzServer.Hertz, httpHandler *handler.UserHTTPHandler, authMiddleware *middleware.AuthMiddleware) {
	//健康检查
	h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"status":  "ok",
			"time":    time.Now().Unix(),
			"service": "user-service",
		})
	})

	//API文档
	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"service":     "User Service",
			"description": "用户微服务",
			"apis": []utils.H{
				{"method": "POST", "path": "/api/v1/auth/register", "desc": "用户注册", "auth": false},
				{"method": "POST", "path": "/api/v1/auth/login", "desc": "用户登录", "auth": false},
				{"method": "POST", "path": "/api/v1/auth/logout", "desc": "用户登出", "auth": true},
				{"method": "PUT", "path": "/api/v1/user/profile", "desc": "更新用户信息", "auth": true},
				{"method": "PUT", "path": "/api/v1/user/password", "desc": "修改密码", "auth": true},
				{"method": "PUT", "path": "/api/v1/user/email", "desc": "修改邮箱", "auth": true},
				{"method": "PUT", "path": "/api/v1/user/phone", "desc": "修改手机号", "auth": true},
				{"method": "GET", "path": "/api/v1/user/:id", "desc": "获取用户资料", "auth": true},
				{"method": "GET", "path": "/api/v1/user/:id/status", "desc": "获取用户状态", "auth": true},
				{"method": "GET", "path": "/api/v1/admin/users", "desc": "管理员：用户列表", "auth": true, "admin": true},
				{"method": "POST", "path": "/api/v1/admin/users/:id/ban", "desc": "管理员：封禁用户", "auth": true, "admin": true},
				{"method": "POST", "path": "/api/v1/admin/users/:id/unban", "desc": "管理员：解封用户", "auth": true, "admin": true},
				{"method": "DELETE", "path": "/api/v1/admin/users/:id", "desc": "管理员：删除用户", "auth": true, "admin": true},
				{"method": "POST", "path": "/api/v1/admin/users/:id/restore", "desc": "管理员：恢复用户", "auth": true, "admin": true},
			},
		})
	})

	//公开路由（不需要认证）
	h.POST("/api/v1/auth/register", httpHandler.Register)
	h.POST("/api/v1/auth/login", httpHandler.Login)

	//需要认证的路由
	authGroup := h.Group("/api/v1")
	authGroup.Use(authMiddleware.HertzMiddleware())
	{
		//用户操作
		authGroup.PUT("/user/profile", httpHandler.UpdateUser)
		authGroup.PUT("/user/password", httpHandler.ChangePassword)
		authGroup.PUT("/user/email", httpHandler.ChangeEmail)
		authGroup.PUT("/user/phone", httpHandler.ChangePhone)
		authGroup.GET("/user/:id", httpHandler.GetUserProfile)
		authGroup.GET("/user/:id/status", httpHandler.GetUserStatus)
		authGroup.POST("/auth/logout", httpHandler.Logout)
	}

	//管理员路由（需要管理员权限）
	adminGroup := h.Group("/api/v1/admin")
	adminGroup.Use(authMiddleware.AdminMiddleware())
	{
		adminGroup.POST("/users/:id/ban", httpHandler.BanUser)
		adminGroup.POST("/users/:id/unban", httpHandler.UnbanUser)
		adminGroup.DELETE("/users/:id", httpHandler.DeleteUser)
		adminGroup.POST("/users/:id/restore", httpHandler.RestoreUser)
		adminGroup.PUT("/users/:id/status", httpHandler.UpdateUserStatus)
		adminGroup.GET("/users", httpHandler.ListUsers)
		adminGroup.GET("/users/search", httpHandler.SearchUsers)
		adminGroup.GET("/users/count", httpHandler.CountUsers)
		adminGroup.GET("/users/count-by-status", httpHandler.CountByStatus)
	}
}

// 关闭
func gracefulShutdown(httpServer *hertzServer.Hertz, kitexServer kitexServer.Server, db *gorm.DB) {
	//创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("💖开始关闭...💖")

	//关闭HTTP服务器
	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("😭关闭HTTP服务器失败: %v😭", err)
		} else {
			log.Println("💖HTTP服务器已关闭💖")
		}
	}

	//关闭Kitex服务器
	if kitexServer != nil {
		if err := kitexServer.Stop(); err != nil {
			log.Printf("😭关闭Kitex服务器失败: %v😭", err)
		} else {
			log.Println("💖Kitex服务器已关闭💖")
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

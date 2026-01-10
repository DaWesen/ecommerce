package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce/gateway/config"
	"ecommerce/gateway/internal/client"
	"ecommerce/gateway/internal/middleware"
	router "ecommerce/gateway/internal/route"
	"ecommerce/gateway/pkg/logger"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func main() {
	//加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("💥 加载配置失败: %v", err)
	}

	//初始化日志
	logger.InitLogger(cfg)

	//初始化 RPC 客户端
	clientManager, err := client.NewClientManager(cfg)
	if err != nil {
		hlog.Fatalf("💥 初始化客户端失败: %v", err)
	}
	defer clientManager.Close()

	hlog.Info("✅ RPC 客户端初始化成功")

	//创建 Hertz 服务器
	h := server.New(
		server.WithHostPorts(fmt.Sprintf(":%d", cfg.Server.Port)),
		server.WithMaxRequestBodySize(10*1024*1024), // 10MB
		server.WithReadTimeout(30*time.Second),
		server.WithWriteTimeout(30*time.Second),
		server.WithIdleTimeout(120*time.Second),
	)

	//注册全局中间件
	registerGlobalMiddleware(h, cfg)

	//注册路由
	router.RegisterRoutes(h, clientManager, cfg)

	//启动服务器
	go func() {
		hlog.Infof("🚀 网关服务启动成功，监听端口: %d", cfg.Server.Port)
		hlog.Infof("🌐 访问地址: http://localhost:%d", cfg.Server.Port)
		hlog.Infof("📊 环境: %s", cfg.Server.Env)
		hlog.Infof("📝 日志级别: %s", cfg.Log.Level)

		if err := h.Run(); err != nil {
			hlog.Fatalf("💥 服务器启动失败: %v", err)
		}
	}()

	//优雅关闭
	waitForShutdown(h)
}

// registerGlobalMiddleware 注册全局中间件
func registerGlobalMiddleware(h *server.Hertz, cfg *config.Config) {
	// CORS 跨域
	h.Use(middleware.CORS())

	// 请求日志
	h.Use(middleware.RequestLogger())

	// 异常恢复
	h.Use(middleware.Recovery())

	// 限流（如果启用）
	if cfg.RateLimit.Enable {
		h.Use(middleware.RateLimiter(cfg.RateLimit))
	}

}

// waitForShutdown 等待关闭信号
func waitForShutdown(h *server.Hertz) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	hlog.Info("⏳ 收到关闭信号，开始关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		hlog.Errorf("关闭服务器失败: %v", err)
	} else {
		hlog.Info("✅ 服务器已关闭")
	}
}

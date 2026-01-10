package router

import (
	"context"
	"ecommerce/gateway/config"
	"ecommerce/gateway/internal/client"
	"ecommerce/gateway/internal/handler"
	"ecommerce/gateway/internal/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

func RegisterRoutes(h *server.Hertz, clientManager *client.ClientManager, cfg *config.Config) {
	// API 分组
	api := h.Group("/api")

	// v1 版本
	v1 := api.Group("/v1")

	// 公共路由（不需要认证）
	public := v1.Group("")
	registerPublicRoutes(public, clientManager)

	// 受保护路由（需要认证）
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(cfg.JWT))
	registerProtectedRoutes(protected, clientManager)

	// 管理路由（需要管理员权限）
	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuth(cfg.JWT))
	admin.Use(middleware.AdminOnly())
	registerAdminRoutes(admin, clientManager)
}

func registerPublicRoutes(group *route.RouterGroup, clientManager *client.ClientManager) {
	// 健康检查
	group.GET("/health", handler.HealthCheck)
	group.GET("/status", handler.ServiceStatus)

	// 用户相关（公共接口）
	group.POST("/users/register", handler.CreateUser(clientManager))
	group.POST("/users/login", handler.UserLogin(clientManager))
	group.GET("/users/:id", handler.GetUserProfile(clientManager))

	// 商品相关（公共接口）
	group.GET("/products", handler.ListProducts(clientManager))
	group.GET("/products/:id", handler.GetProduct(clientManager))
	group.GET("/products/category/:category", handler.GetProductsByCategory(clientManager))
	group.POST("/products/search", handler.SearchProducts(clientManager))

	// 订单相关（部分公共接口）
	group.GET("/orders/:order_no", handler.GetOrder(clientManager))
}

func registerProtectedRoutes(group *route.RouterGroup, clientManager *client.ClientManager) {
	// 用户相关
	group.PUT("/users/profile", handler.UpdateUserProfile(clientManager))
	group.PUT("/users/password", handler.ChangePassword(clientManager))
	group.GET("/users/me", handler.GetCurrentUser(clientManager))
	group.POST("/users/logout", handler.UserLogout(clientManager))

	// 订单相关
	group.POST("/orders", handler.CreateOrder(clientManager))
	group.GET("/orders", handler.ListOrders(clientManager))
	group.POST("/orders/:order_no/pay", handler.PayOrder(clientManager))
	group.POST("/orders/:order_no/cancel", handler.CancelOrder(clientManager))
	group.POST("/orders/:order_no/receive", handler.ConfirmReceipt(clientManager))
	group.POST("/orders/:order_no/refund", handler.ApplyRefund(clientManager))
}

func registerAdminRoutes(group *route.RouterGroup, clientManager *client.ClientManager) {
	// 用户管理
	group.GET("/users", handler.ListUsers(clientManager))
	group.POST("/users/search", handler.SearchUsers(clientManager))
	group.PUT("/users/:id/status", handler.UpdateUserStatus(clientManager))
	group.DELETE("/users/:id", handler.DeleteUser(clientManager))

	// 商品管理
	group.POST("/products", handler.CreateProduct(clientManager))
	group.PUT("/products/:id", handler.UpdateProduct(clientManager))
	group.DELETE("/products/:id", handler.DeleteProduct(clientManager))
	group.POST("/products/:id/online", handler.OnlineProduct(clientManager))
	group.POST("/products/:id/offline", handler.OfflineProduct(clientManager))
	group.POST("/products/search", handler.AdminSearchProducts(clientManager))

	// 订单管理
	group.GET("/orders/all", handler.ListAllOrders(clientManager))
	group.POST("/orders/:order_no/ship", handler.ShipOrder(clientManager))
	group.POST("/orders/refunds/:refund_no/process", handler.ProcessRefund(clientManager))
	group.GET("/stats/orders", handler.GetOrderStats(clientManager))
}

func registerDocRoutes(h *server.Hertz) {
	// API 文档
	h.GET("/docs", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(200, map[string]interface{}{
			"service": "E-Commerce Gateway API",
			"version": "1.0.0",
			"docs": map[string]string{
				"swagger": "/swagger/index.html",
				"openapi": "/openapi.json",
			},
		})
	})

	// Swagger UI
	h.Static("/swagger", "./docs/swagger")

	// OpenAPI 规范
	h.GET("/openapi.json", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(200, map[string]interface{}{
			"openapi": "3.0.0",
			"info": map[string]interface{}{
				"title":       "E-Commerce Gateway API",
				"version":     "1.0.0",
				"description": "电商系统API网关",
			},
			"paths": map[string]interface{}{
				"/api/v1/health": map[string]interface{}{
					"get": map[string]interface{}{
						"summary": "健康检查",
						"responses": map[string]interface{}{
							"200": map[string]interface{}{
								"description": "服务健康状态",
							},
						},
					},
				},
			},
		})
	})
}

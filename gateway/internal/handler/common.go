package handler

import (
	"context"
	"ecommerce/gateway/pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
)

// HealthCheck 健康检查
func HealthCheck(c context.Context, ctx *app.RequestContext) {
	response.Success(ctx, map[string]interface{}{
		"status":  "ok",
		"service": "ecommerce-gateway",
		"version": "1.0.0",
	})
}

// ServiceStatus 服务状态
func ServiceStatus(c context.Context, ctx *app.RequestContext) {
	response.Success(ctx, map[string]interface{}{
		"status":  "running",
		"service": "ecommerce-gateway",
		"version": "1.0.0",
		"apis": map[string]interface{}{
			"public": []string{
				"/api/v1/health",
				"/api/v1/status",
				"/api/v1/users/register",
				"/api/v1/users/login",
				"/api/v1/products",
				"/api/v1/products/:id",
				"/api/v1/products/category/:category",
				"/api/v1/products/search",
			},
			"protected": []string{
				"/api/v1/users/profile",
				"/api/v1/users/password",
				"/api/v1/users/me",
				"/api/v1/orders",
				"/api/v1/orders/:order_no",
				"/api/v1/orders/:order_no/pay",
				"/api/v1/orders/:order_no/cancel",
				"/api/v1/orders/:order_no/refund",
			},
			"admin": []string{
				"/api/v1/admin/users",
				"/api/v1/admin/users/:id/status",
				"/api/v1/admin/users/:id",
				"/api/v1/admin/products",
				"/api/v1/admin/products/:id",
				"/api/v1/admin/orders/all",
				"/api/v1/admin/orders/:order_no/ship",
				"/api/v1/admin/stats/orders",
			},
		},
	})
}

package middleware

import (
	"context"
	"ecommerce/gateway/config"

	"github.com/cloudwego/hertz/pkg/app"
	"golang.org/x/time/rate"
)

// 简单的IP限流器
var ipLimiters = make(map[string]*rate.Limiter)

// RateLimiter 限流中间件
func RateLimiter(cfg config.RateLimitConfig) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		if !cfg.Enable {
			ctx.Next(c)
			return
		}

		// 获取客户端IP
		clientIP := ctx.ClientIP()

		// 获取或创建限流器
		limiter, exists := ipLimiters[clientIP]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(cfg.Requests), cfg.Burst)
			ipLimiters[clientIP] = limiter
		}

		// 尝试获取令牌
		if !limiter.Allow() {
			ctx.JSON(429, map[string]interface{}{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			ctx.Abort()
			return
		}

		ctx.Next(c)
	}
}

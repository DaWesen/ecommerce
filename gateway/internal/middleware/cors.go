package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CORS 跨域中间件
func CORS() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Accept-Language, Cache-Control, X-Requested-With")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Max-Age", "86400") // 24小时

		if string(ctx.Request.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next(c)
	}
}

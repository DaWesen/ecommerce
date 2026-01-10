package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// RequestLogger 请求日志中间件
func RequestLogger() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		ctx.Next(c)

		latency := time.Since(start)
		hlog.CtxInfof(c, "Request | %s %s | %d | %v | %s | %s",
			ctx.Request.Method(),
			ctx.Request.URI().Path(),
			ctx.Response.StatusCode(),
			latency,
			ctx.ClientIP(),
			string(ctx.Request.Header.UserAgent()),
		)
	}
}

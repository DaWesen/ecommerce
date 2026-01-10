package middleware

import (
	"context"
	"runtime/debug"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Recovery 异常恢复中间件
func Recovery() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				hlog.CtxErrorf(c, "Recovery | panic recovered: %v\n%s", err, debug.Stack())
				ctx.JSON(500, map[string]interface{}{
					"code":    500,
					"message": "服务器内部错误",
				})
				ctx.Abort()
			}
		}()

		ctx.Next(c)
	}
}

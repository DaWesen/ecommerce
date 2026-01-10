package response

import (
	"github.com/cloudwego/hertz/pkg/app"
)

// 成功响应
func Success(ctx *app.RequestContext, data interface{}) {
	ctx.JSON(200, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    data,
		"success": true,
	})
}

// 错误响应
func Error(ctx *app.RequestContext, code int, message string) {
	ctx.JSON(code, map[string]interface{}{
		"code":    code,
		"message": message,
		"success": false,
	})
}

// 带分页的成功响应
func SuccessWithPagination(ctx *app.RequestContext, data interface{}, total int64, page, pageSize int) {
	ctx.JSON(200, map[string]interface{}{
		"code":      0,
		"message":   "success",
		"data":      data,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"success":   true,
	})
}

// 自定义响应
func Custom(ctx *app.RequestContext, code int, data map[string]interface{}) {
	ctx.JSON(code, data)
}

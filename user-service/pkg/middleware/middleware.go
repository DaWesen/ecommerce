package middleware

import (
	"context"
	"ecommerce/user-service/pkg/jwt"
	"errors"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/endpoint"
)

var (
	ErrTokenMissing = errors.New("token is missing")
	ErrTokenInvalid = errors.New("token is invalid")
	ErrNoPermission = errors.New("no permission")
)

// JWT认证中间件
type AuthMiddleware struct {
	jwtManager    *jwt.JWTManager
	excludedPaths []string
}

// 创建认证中间件
func NewAuthMiddleware(jwtManager *jwt.JWTManager, excludedPaths []string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager:    jwtManager,
		excludedPaths: excludedPaths,
	}
}

// Hertz HTTP中间件
func (m *AuthMiddleware) HertzMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		//检查是否在排除列表中
		path := string(c.Request.URI().Path())
		for _, excludedPath := range m.excludedPaths {
			if strings.HasPrefix(path, excludedPath) {
				c.Next(ctx)
				return
			}
		}
		//从请求头获取token
		token := extractTokenFromHeader(c)
		if token == "" {
			c.JSON(401, map[string]interface{}{
				"code":    401,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}
		//验证token
		claims, err := m.jwtManager.VerifyAccessToken(token)
		if err != nil {
			hlog.Warnf("Token验证失败: %v", err)
			c.JSON(401, map[string]interface{}{
				"code":    401,
				"message": "认证令牌无效或已过期",
			})
			c.Abort()
			return
		}
		//将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_status", claims.Status)
		c.Set("is_admin", claims.IsAdmin)
		//继续处理
		c.Next(ctx)
	}
}

// Kitex RPC中间件
func (m *AuthMiddleware) KitexMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		//从metainfo中获取token
		authHeader, exists := metainfo.GetValue(ctx, "authorization")
		if !exists {
			return ErrTokenMissing
		}

		//提取Bearer token
		token := extractTokenFromString(authHeader)
		if token == "" {
			return ErrTokenInvalid
		}
		//验证token
		claims, err := m.jwtManager.VerifyAccessToken(token)
		if err != nil {
			return err
		}
		//将用户信息存入上下文，供后续处理使用
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_status", claims.Status)
		ctx = context.WithValue(ctx, "is_admin", claims.IsAdmin)
		// 添加日志记录
		hlog.CtxInfof(ctx, "用户认证成功: user_id=%d, email=%s, is_admin=%v",
			claims.UserID, claims.Email, claims.IsAdmin)

		return next(ctx, req, resp)
	}
}

// Kitex RPC管理员权限中间件
func (m *AuthMiddleware) KitexAdminMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		//从metainfo中获取token
		authHeader, exists := metainfo.GetValue(ctx, "authorization")
		if !exists {
			return ErrTokenMissing
		}

		//提取Bearer token
		token := extractTokenFromString(authHeader)
		if token == "" {
			return ErrTokenInvalid
		}

		//验证token
		claims, err := m.jwtManager.VerifyAccessToken(token)
		if err != nil {
			return err
		}

		//检查是否是管理员
		if !claims.IsAdmin {
			return ErrNoPermission
		}

		//将用户信息存入上下文
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_status", claims.Status)
		ctx = context.WithValue(ctx, "is_admin", claims.IsAdmin)
		hlog.CtxInfof(ctx, "管理员认证成功: user_id=%d, email=%s", claims.UserID, claims.Email)
		return next(ctx, req, resp)
	}
}

// 管理员权限中间件
func (m *AuthMiddleware) AdminMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		//先执行认证
		m.HertzMiddleware()(ctx, c)
		if c.IsAborted() {
			return
		}
		//检查是否是管理员
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(403, map[string]interface{}{
				"code":    403,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}

// 从请求头提取token
func extractTokenFromHeader(c *app.RequestContext) string {
	//从Authorization头获取
	authHeaderBytes := c.GetHeader("Authorization")
	if len(authHeaderBytes) == 0 {
		return ""
	}
	authHeader := string(authHeaderBytes)
	//检查Bearer token格式
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// 从字符串提取
func extractTokenFromString(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	//检查Bearer token格式
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

// 从Hertz上下文中获取用户ID
func GetUserIDFromContext(c *app.RequestContext) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int64), true
}

// 从Hertz上下文中获取用户邮箱
func GetUserEmailFromContext(c *app.RequestContext) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	return email.(string), true
}

// 从Hertz上下文中判断是否是管理员
func IsAdminFromContext(c *app.RequestContext) bool {
	isAdmin, exists := c.Get("is_admin")
	if !exists {
		return false
	}
	return isAdmin.(bool)
}

// 从Hertz上下文中获取用户状态
func GetUserStatusFromContext(c *app.RequestContext) (string, bool) {
	status, exists := c.Get("user_status")
	if !exists {
		return "", false
	}
	return status.(string), true
}

// 从Hertz上下文中获取所有用户信息
func GetUserInfoFromContext(c *app.RequestContext) (map[string]interface{}, bool) {
	userInfo := make(map[string]interface{})
	if userID, exists := c.Get("user_id"); exists {
		userInfo["user_id"] = userID.(int64)
	}
	if email, exists := c.Get("user_email"); exists {
		userInfo["user_email"] = email.(string)
	}
	if status, exists := c.Get("user_status"); exists {
		userInfo["user_status"] = status.(string)
	}
	if isAdmin, exists := c.Get("is_admin"); exists {
		userInfo["is_admin"] = isAdmin.(bool)
	}
	return userInfo, len(userInfo) > 0
}

// 从Kitex上下文中获取用户ID
func GetUserIDFromKitexContext(ctx context.Context) (int64, bool) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return 0, false
	}
	id, ok := userID.(int64)
	if !ok {
		return 0, false
	}
	return id, true
}

// 从Kitex上下文中获取用户邮箱
func GetUserEmailFromKitexContext(ctx context.Context) (string, bool) {
	email := ctx.Value("user_email")
	if email == nil {
		return "", false
	}
	e, ok := email.(string)
	if !ok {
		return "", false
	}
	return e, true
}

// 从Kitex上下文中判断是否是管理员
func IsAdminFromKitexContext(ctx context.Context) bool {
	isAdmin := ctx.Value("is_admin")
	if isAdmin == nil {
		return false
	}
	admin, ok := isAdmin.(bool)
	if !ok {
		return false
	}
	return admin
}

// 从Kitex上下文中获取用户状态
func GetUserStatusFromKitexContext(ctx context.Context) (string, bool) {
	status := ctx.Value("user_status")
	if status == nil {
		return "", false
	}
	s, ok := status.(string)
	if !ok {
		return "", false
	}
	return s, true
}

package middleware

import (
	"context"
	"ecommerce/user-service/kitex_gen/api"
	"ecommerce/user-service/pkg/jwt"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
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
		//获取当前调用的方法名
		methodName := ""
		if ri := rpcinfo.GetRPCInfo(ctx); ri != nil {
			methodName = ri.Invocation().MethodName()
		}

		log.Printf("KitexMiddleware: 调用方法: %s", methodName)

		//定义不需要认证的公开方法
		publicMethods := map[string]bool{
			"Register": true,
			"Login":    true,
		}

		//如果是公开方法，直接跳过认证
		if publicMethods[methodName] {
			log.Printf("KitexMiddleware: %s 是公开方法，跳过认证", methodName)
			return next(ctx, req, resp)
		}

		//尝试从metainfo获取token
		var token string
		possibleHeaders := []string{
			"authorization",
			"Authorization",
			"AUTHORIZATION",
			"x-auth-token",
			"X-Auth-Token",
			"token",
			"Token",
			"x-token",
			"X-Token",
		}

		for _, header := range possibleHeaders {
			if val, exists := metainfo.GetValue(ctx, header); exists {
				token = val
				log.Printf("KitexMiddleware: 从metainfo[%s]获取到token", header)
				break
			}
		}

		//如果metainfo中没有，尝试从请求结构体中获取
		if token == "" {
			token = extractTokenFromArgs(req)
			if token != "" {
				log.Printf("KitexMiddleware: 从Args结构体中提取到token")
			}
		}

		if token == "" {
			log.Printf("KitexMiddleware: 没有找到任何认证信息")
			return ErrTokenMissing
		}

		//提取Bearer token（如果包含Bearer前缀）
		token = extractTokenFromString(token)
		if token == "" {
			log.Printf("KitexMiddleware: token为空或格式不正确")
			return ErrTokenInvalid
		}

		//验证token
		claims, err := m.jwtManager.VerifyAccessToken(token)
		if err != nil {
			log.Printf("KitexMiddleware: token验证失败: %v", err)
			return err
		}

		log.Printf("KitexMiddleware: 用户认证成功: user_id=%d, email=%s, is_admin=%v",
			claims.UserID, claims.Email, claims.IsAdmin)

		//将用户信息存入上下文，供后续处理使用
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_status", claims.Status)
		ctx = context.WithValue(ctx, "is_admin", claims.IsAdmin)

		return next(ctx, req, resp)
	}
}

// 添加辅助函数
func safeSubstring(s string, start, end int) string {
	if s == "" {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return ""
	}
	return s[start:end]
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
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	//如果没有Bearer前缀，直接返回
	return authHeader
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

// 从Args结构体中提取请求并获取token
func extractTokenFromArgs(req interface{}) string {
	//使用反射来访问Args结构体的Req字段
	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() == reflect.Ptr {
		reqValue = reqValue.Elem()
	}

	//Args结构体通常有一个名为"Req"的字段
	reqField := reqValue.FieldByName("Req")
	if !reqField.IsValid() || reqField.IsZero() {
		log.Printf("extractTokenFromArgs: 无法找到Req字段")
		return ""
	}

	//获取实际的请求对象
	actualReq := reqField.Interface()

	//根据实际请求类型提取token
	switch r := actualReq.(type) {
	case *api.GetUserProfileReq:
		return r.Token
	case *api.UpdateUserReq:
		return r.Token
	case *api.ChangePasswordReq:
		return r.Token
	case *api.ChangeEmailReq:
		return r.Token
	case *api.ChangePhoneReq:
		return r.Token
	case *api.LogoutReq:
		return r.Token
	case *api.GetUserStatusReq:
		return r.Token
	case *api.BanUserReq:
		return r.Token
	case *api.UnbanUserReq:
		return r.Token
	case *api.DeleteUserReq:
		return r.Token
	case *api.RestoreUserReq:
		return r.Token
	case *api.UpdateUserStatusReq:
		return r.Token
	case *api.ListUsersReq:
		return r.Token
	case *api.SearchUsersReq:
		return r.Token
	case *api.CountUsersReq:
		return r.Token
	case *api.CountByStatusReq:
		return r.Token
	case *api.UpdatePasswordReq:
		return r.Token
	case *api.UpdateEmailReq:
		return r.Token
	case *api.UpdatePhoneReq:
		return r.Token
	case *api.UpdateUserProfileReq:
		return r.Token
	default:
		log.Printf("extractTokenFromArgs: 未知的请求类型: %T", actualReq)
		return ""
	}
}

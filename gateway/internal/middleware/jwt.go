package middleware

import (
	"context"
	"strings"
	"time"

	"ecommerce/gateway/config"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/golang-jwt/jwt/v5"
)

// JWT中间件
func JWTAuth(cfg config.JWTConfig) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 从请求头获取token
		authHeader := string(ctx.GetHeader("Authorization"))
		if authHeader == "" {
			ctx.JSON(401, map[string]interface{}{
				"code":    401,
				"message": "未提供认证令牌",
			})
			ctx.Abort()
			return
		}

		// 检查Bearer格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(401, map[string]interface{}{
				"code":    401,
				"message": "认证令牌格式错误",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// 解析和验证token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		})

		if err != nil {
			hlog.Errorf("JWT解析失败: %v", err)
			ctx.JSON(401, map[string]interface{}{
				"code":    401,
				"message": "认证令牌无效",
			})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 检查token是否过期
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					ctx.JSON(401, map[string]interface{}{
						"code":    401,
						"message": "认证令牌已过期",
					})
					ctx.Abort()
					return
				}
			}

			// 验证签发者
			if iss, ok := claims["iss"].(string); ok && iss != cfg.Issuer {
				ctx.JSON(401, map[string]interface{}{
					"code":    401,
					"message": "认证令牌签发者无效",
				})
				ctx.Abort()
				return
			}

			// 将用户信息存储到上下文中
			if userID, ok := claims["user_id"].(float64); ok {
				ctx.Set("user_id", int64(userID))
			}
			if username, ok := claims["username"].(string); ok {
				ctx.Set("username", username)
			}
			if role, ok := claims["role"].(string); ok {
				ctx.Set("role", role)
			}

			ctx.Next(c)
		} else {
			ctx.JSON(401, map[string]interface{}{
				"code":    401,
				"message": "认证令牌验证失败",
			})
			ctx.Abort()
		}
	}
}

// AdminOnly 管理员权限中间件
func AdminOnly() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		role, exists := ctx.Get("role")
		if !exists || role != "admin" {
			ctx.JSON(403, map[string]interface{}{
				"code":    403,
				"message": "需要管理员权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next(c)
	}
}

// GenerateJWTToken 生成JWT令牌
func GenerateJWTToken(cfg config.JWTConfig, userID int64, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"iss":      cfg.Issuer,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * time.Duration(cfg.ExpireHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

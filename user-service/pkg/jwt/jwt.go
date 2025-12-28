package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT配置
type Config struct {
	SecretKey     string        `json:"secret_key"`     //密钥
	Issuer        string        `json:"issuer"`         //签发者
	AccessExpire  time.Duration `json:"access_expire"`  //访问令牌过期时间
	RefreshExpire time.Duration `json:"refresh_expire"` //刷新令牌过期时间
	Algorithm     string        `json:"algorithm"`      //算法
}

// 默认配置
var DefaultConfig = Config{
	SecretKey:     "your-secret-key-change-in-production",
	Issuer:        "ecommerce-user-service",
	AccessExpire:  24 * time.Hour,
	RefreshExpire: 7 * 24 * time.Hour,
	Algorithm:     "HS256",
}

// 自定义Claims
type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Status   int32  `json:"status"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// JWT管理器
type JWTManager struct {
	config Config
}

// 创建JWT管理器
func NewJWTManager(config Config) *JWTManager {
	return &JWTManager{
		config: config,
	}
}

// 生成访问令牌
func (m *JWTManager) GenerateAccessToken(userID int64, email, username string, status int32, isAdmin bool) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Status:   status,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.AccessExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    m.config.Issuer,
			Subject:   "access_token",
			ID:        fmt.Sprintf("%d-%d", userID, time.Now().UnixNano()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SecretKey))
}

// 生成刷新令牌
func (m *JWTManager) GenerateRefreshToken(userID int64) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    m.config.Issuer,
			Subject:   "refresh_token",
			ID:        fmt.Sprintf("refresh-%d-%d", userID, time.Now().UnixNano()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SecretKey))
}

// 生成令牌对
func (m *JWTManager) GenerateTokenPair(userID int64, email, username string, status int32, isAdmin bool) (accessToken, refreshToken string, err error) {
	accessToken, err = m.GenerateAccessToken(userID, email, username, status, isAdmin)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = m.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// 验证令牌并返回Claims
func (m *JWTManager) VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 验证访问令牌
func (m *JWTManager) VerifyAccessToken(tokenString string) (*CustomClaims, error) {
	claims, err := m.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查是否为访问令牌
	if claims.Subject != "access_token" {
		return nil, errors.New("not an access token")
	}

	return claims, nil
}

// 验证刷新令牌
func (m *JWTManager) VerifyRefreshToken(tokenString string) (*CustomClaims, error) {
	claims, err := m.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	//检查是否为刷新令牌
	if claims.Subject != "refresh_token" {
		return nil, errors.New("not a refresh token")
	}
	return claims, nil
}

// 刷新令牌
func (m *JWTManager) RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	claims, err := m.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	//生成新的令牌对
	newAccessToken, newRefreshToken, err = m.GenerateTokenPair(
		claims.UserID,
		claims.Email,
		claims.Username,
		claims.Status,
		claims.IsAdmin,
	)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

// 从令牌中获取用户ID
func (m *JWTManager) GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := m.VerifyToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// 检查令牌是否有效
func (m *JWTManager) IsTokenValid(tokenString string) bool {
	_, err := m.VerifyToken(tokenString)
	return err == nil
}

// 获取令牌过期时间
func (m *JWTManager) GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := m.VerifyToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return claims.ExpiresAt.Time, nil
}

// 获取令牌签发时间
func (m *JWTManager) GetTokenIssuedAt(tokenString string) (time.Time, error) {
	claims, err := m.VerifyToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return claims.IssuedAt.Time, nil
}

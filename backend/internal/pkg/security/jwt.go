// Package security 提供 JWT 生成/解析和密码哈希工具。
package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go_web/internal/config"
)

// Claims JWT 载荷，包含用户身份信息和标准注册声明。
type Claims struct {
	UserID   uint     `json:"user_id"`   // 用户 ID
	Username string   `json:"username"`  // 用户名
	Roles    []string `json:"roles"`     // 用户角色编码列表
	jwt.RegisteredClaims
}

// JWTManager 管理 JWT Token 的生成与解析。
type JWTManager struct {
	cfg config.JWTConfig
}

// NewJWTManager 创建 JWT 管理器。
func NewJWTManager(cfg config.JWTConfig) *JWTManager {
	return &JWTManager{cfg: cfg}
}

// Generate 为用户生成 JWT Token，内含用户 ID、用户名和角色列表。
func (m *JWTManager) Generate(userID uint, username string, roles []string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.cfg.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.cfg.AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.Secret))
}

// Parse 解析并验证 JWT Token，返回 Claims。
func (m *JWTManager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(m.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

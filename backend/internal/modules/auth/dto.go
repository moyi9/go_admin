// Package auth 提供用户认证相关的登录与当前用户查询功能。
package auth

import "go_web/internal/modules/rbac"

// LoginRequest 登录请求体。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录成功响应，包含 JWT Token 和用户信息。
type LoginResponse struct {
	AccessToken string    `json:"access_token"` // JWT Token 字符串
	TokenType   string    `json:"token_type"`   // Token 类型，固定 Bearer
	User        rbac.User `json:"user"`         // 当前登录用户信息
}

// UpdatePasswordRequest 修改密码请求体。
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// UpdateProfileRequest 更新个人信息请求体。
type UpdateProfileRequest struct {
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
	AvatarURL string `json:"avatar_url"`
}

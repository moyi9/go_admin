package auth

import (
	"github.com/gin-gonic/gin"
	"go_web/internal/modules/audit"
	"go_web/internal/modules/notification"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/response"
)

// Handler 处理认证相关 HTTP 请求。
type Handler struct {
	service *Service
	audit   *audit.Service
	notif   *notification.Service
}

// NewHandler 创建认证 Handler。
func NewHandler(service *Service, auditService *audit.Service, notifService *notification.Service) *Handler {
	return &Handler{service: service, audit: auditService, notif: notifService}
}

// Login 用户登录。验证用户名密码，成功后返回 JWT Token。
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, apperror.CodeInvalidArgument, err.Error())
		return
	}

	result, err := h.service.Login(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	// 记录登录日志
	if h.audit != nil {
		h.audit.Log(&audit.AuditLog{
			UserID:    result.User.ID,
			Username:  result.User.Username,
			Action:    audit.ActionLogin,
			Resource:  audit.ResourceUser,
			Detail:    "用户登录",
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		})
	}
	response.OK(c, result)
}

// Me 获取当前登录用户信息（从 JWT 中提取 user_id）。
func (h *Handler) Me(c *gin.Context) {
	v, ok := c.Get("user_id")
	if !ok {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}
	userID, ok := v.(uint)
	if !ok {
		response.Error(c, 401, apperror.CodeUnauthorized, "invalid user id")
		return
	}
	user, err := h.service.Me(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, user)
}

// UpdatePassword 修改当前登录用户密码。
func (h *Handler) UpdatePassword(c *gin.Context) {
	v, ok := c.Get("user_id")
	if !ok {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}
	userID, ok := v.(uint)
	if !ok {
		response.Error(c, 401, apperror.CodeUnauthorized, "invalid user id")
		return
	}
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, apperror.CodeInvalidArgument, err.Error())
		return
	}
	if err := h.service.UpdatePassword(userID, req); err != nil {
		response.FromError(c, err)
		return
	}
	if h.audit != nil {
		username, _ := c.Get("username")
		uname, _ := username.(string)
		h.audit.Log(&audit.AuditLog{
			UserID:    userID,
			Username:  uname,
			Action:    audit.ActionUpdate,
			Resource:  audit.ResourceUser,
			Detail:    "修改密码",
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		})
	}
	if h.notif != nil {
			username, _ := c.Get("username")
			uname, _ := username.(string)
			h.notif.Send(userID, notification.TypeSecurity, "密码修改", "您的登录密码已修改成功", "")
			_ = uname
		}
		response.OK(c, gin.H{"message": "password updated"})
}

// UpdateProfile 更新当前登录用户个人信息。
func (h *Handler) UpdateProfile(c *gin.Context) {
	v, ok := c.Get("user_id")
	if !ok {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}
	userID, ok := v.(uint)
	if !ok {
		response.Error(c, 401, apperror.CodeUnauthorized, "invalid user id")
		return
	}
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, apperror.CodeInvalidArgument, err.Error())
		return
	}
	user, err := h.service.UpdateProfile(userID, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	if h.audit != nil {
		username, _ := c.Get("username")
		uname, _ := username.(string)
		h.audit.Log(&audit.AuditLog{
			UserID:    userID,
			Username:  uname,
			Action:    audit.ActionUpdate,
			Resource:  audit.ResourceUser,
			Detail:    "更新个人资料",
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		})
	}
	if h.notif != nil {
			username, _ := c.Get("username")
			uname, _ := username.(string)
			h.notif.Send(userID, notification.TypeSystem, "资料更新", "您的个人信息已更新", "")
			_ = uname
		}
		response.OK(c, user)
}

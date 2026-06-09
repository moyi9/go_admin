package notification

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go_web/internal/modules/rbac"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/response"
	"gorm.io/gorm"
)

// Handler 处理通知相关 HTTP 请求。
type Handler struct {
	service *Service
	db      *gorm.DB
}

// NewHandler 创建通知 Handler。
func NewHandler(service *Service, db *gorm.DB) *Handler {
	return &Handler{service: service, db: db}
}

// SendNotificationRequest 发送通知的请求体。
type SendNotificationRequest struct {
	Type         string `json:"type" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Content      string `json:"content" binding:"required"`
	TargetUserID uint   `json:"target_user_id"` // 0 = 全体用户
}

// ListNotifications 分页查询当前用户的通知列表。
func (h *Handler) ListNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}
	offset, limit := response.Page(c)
	result, err := h.service.List(userID, offset, limit)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.PaginatedOK(c, result.List, result.Total, offset, limit)
}

// CountUnread 返回当前用户的未读通知数。
func (h *Handler) CountUnread(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}
	count, err := h.service.CountUnread(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, gin.H{"count": count})
}

// SendNotification 管理员发送通知给全体或指定用户。
func (h *Handler) SendNotification(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, apperror.CodeInvalidArgument, err.Error())
		return
	}
	senderID := c.GetUint("user_id")
	if senderID == 0 {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}

	if req.TargetUserID > 0 {
		// 发送给指定用户
		if err := h.service.Send(req.TargetUserID, req.Type, req.Title, req.Content, ""); err != nil {
			response.FromError(c, err)
			return
		}
	} else {
		// 发送给全体活跃用户
		var users []rbac.User
		if err := h.db.Where("status = ?", rbac.UserStatusActive).Find(&users).Error; err != nil {
			response.FromError(c, apperror.Wrap(apperror.CodeInternal, "query users failed", err))
			return
		}
		for _, u := range users {
			if err := h.service.Send(u.ID, req.Type, req.Title, req.Content, ""); err != nil {
				response.FromError(c, err)
				return
			}
		}
	}
	response.OK(c, gin.H{"message": "notification sent"})
}

// MarkAsRead 将指定用户的一条通知标记为已读。
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		response.Error(c, 400, apperror.CodeInvalidArgument, "invalid id")
		return
	}
	id := uint(id64)
	if err := h.service.MarkAsRead(userID, id); err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "marked as read"})
}

// MarkAllAsRead 将指定用户的全部通知标记为已读。
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}
	if err := h.service.MarkAllAsRead(userID); err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "all marked as read"})
}

// ListUnreadNotifications 返回当前用户的最新 N 条未读通知（供铃铛下拉预览）。
func (h *Handler) ListUnreadNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, 401, apperror.CodeUnauthorized, "unauthorized")
		return
	}
	list, err := h.service.ListUnread(userID, 5)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, list)
}

package audit

import (
	"github.com/gin-gonic/gin"
	"go_web/internal/pkg/response"
)

// Handler 处理操作日志相关 HTTP 请求。
type Handler struct {
	service *Service
}

// NewHandler 创建操作日志 Handler。
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListAuditLogs 分页查询操作日志列表。
func (h *Handler) ListAuditLogs(c *gin.Context) {
	offset, limit := response.Page(c)
	q := ListQuery{
		Action:   c.Query("action"),
		Resource: c.Query("resource"),
		Keyword:  c.Query("keyword"),
		Since:    c.Query("since"),
		Until:    c.Query("until"),
	}

	result, err := h.service.List(offset, limit, q)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.PaginatedOK(c, result.List, result.Total, offset, limit)
}

package dashboard

import (
	"github.com/gin-gonic/gin"
	"go_web/internal/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Stats 返回仪表盘聚合统计数据。
func (h *Handler) Stats(c *gin.Context) {
	stats, err := h.service.Stats()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, stats)
}

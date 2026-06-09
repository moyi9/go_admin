package response

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Page 从查询参数中提取分页信息，提供合理的默认值和上限。
func Page(c *gin.Context) (offset int, limit int) {
	page := DefaultPage
	pageSize := DefaultPageSize

	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := c.Query("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			pageSize = n
		}
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return (page - 1) * pageSize, pageSize
}

// PaginatedBody 统一的分页响应格式。
type PaginatedBody struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	Total     int64  `json:"total"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	RequestID string `json:"request_id,omitempty"`
}

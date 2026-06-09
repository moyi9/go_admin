// Package response 提供统一 JSON 响应格式，所有 API 响应均通过此包构造。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go_web/internal/pkg/apperror"
)

// Body 标准 JSON 响应体。
type Body struct {
	Code      int    `json:"code"`                // 业务错误码，0 表示成功
	Message   string `json:"message"`             // 提示信息
	Data      any    `json:"data,omitempty"`      // 响应数据
	RequestID string `json:"request_id,omitempty"` // 请求 ID，用于链路追踪
}

// OK 返回 200 成功响应。<｜end▁of▁thinking｜>.OK(c, data)
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:      apperror.CodeOK,
		Message:   "ok",
		Data:      data,
		RequestID: requestID(c),
	})
}

// Created 返回 201 创建成功响应。用于 POST 创建资源后返回新资源。
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Body{
		Code:      apperror.CodeOK,
		Message:   "created",
		Data:      data,
		RequestID: requestID(c),
	})
}

// Error 返回错误响应。status 为 HTTP 状态码，code 为业务错误码，message 为错误描述。
func Error(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Body{
		Code:      code,
		Message:   message,
		RequestID: requestID(c),
	})
}

// PaginatedOK 返回包含分页信息的成功响应。
func PaginatedOK(c *gin.Context, data any, total int64, offset int, limit int) {
	page := 1
	if limit > 0 {
		page = offset/limit + 1
	}
	c.JSON(http.StatusOK, PaginatedBody{
		Code:      apperror.CodeOK,
		Message:   "ok",
		Data:      data,
		Total:     total,
		Page:      page,
		PageSize:  limit,
		RequestID: requestID(c),
	})
}

// FromError 将业务错误转换为 HTTP 响应。优先提取 *apperror.Error，否则返回 500。
func FromError(c *gin.Context, err error) {
	if appErr, ok := apperror.As(err); ok {
		Error(c, statusFromCode(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	Error(c, http.StatusInternalServerError, apperror.CodeInternal, "internal server error")
}

// requestID 从 Gin Context 或请求头中提取 request_id。
func requestID(c *gin.Context) string {
	v, _ := c.Get("request_id")
	if s, ok := v.(string); ok {
		return s
	}
	return c.GetHeader("X-Request-ID")
}

func statusFromCode(code int) int {
	switch code {
	case apperror.CodeInvalidArgument:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperror.CodeForbidden:
		return http.StatusForbidden
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case apperror.CodeConflict:
		return http.StatusConflict
	case apperror.CodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// Package upload 提供通用文件上传功能。
package upload

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go_web/internal/config"
	"go_web/internal/modules/audit"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/response"
)

// Handler 处理文件上传请求。
type Handler struct {
	cfg   config.UploadConfig
	audit *audit.Service
}

// NewHandler 创建上传 Handler。
func NewHandler(cfg config.UploadConfig, auditService *audit.Service) *Handler {
	return &Handler{cfg: cfg, audit: auditService}
}

// Upload 上传文件，返回可访问的 URL。
func (h *Handler) Upload(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	username, _ := c.Get("username")
	uname, _ := username.(string)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, 400, apperror.CodeInvalidArgument, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 校验扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !h.isAllowedExt(ext) {
		response.Error(c, 400, apperror.CodeInvalidArgument, fmt.Sprintf("不支持的文件类型: %s，允许: %s", ext, strings.Join(h.cfg.AllowedExts, ", ")))
		return
	}

	// 校验大小
	if header.Size > h.cfg.MaxSize {
		response.Error(c, 400, apperror.CodeInvalidArgument, fmt.Sprintf("文件过大（最大 %d MB）", h.cfg.MaxSize>>20))
		return
	}

	// 确保上传目录存在
	if err := os.MkdirAll(h.cfg.Dir, 0755); err != nil {
		response.Error(c, 500, apperror.CodeInternal, "创建上传目录失败")
		return
	}

	// 生成唯一文件名：{userID}_{timestamp}_{4位随机}.{ext}
	filename := fmt.Sprintf("%d_%d_%04d%s", uid, time.Now().UnixMilli(), rand.Intn(10000), ext)
	dest := filepath.Join(h.cfg.Dir, filename)

	out, err := os.Create(dest)
	if err != nil {
		response.Error(c, 500, apperror.CodeInternal, "保存文件失败")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		response.Error(c, 500, apperror.CodeInternal, "写入文件失败")
		return
	}

	url := "/uploads/" + filename

	// 记录操作日志
	if h.audit != nil {
		h.audit.Log(&audit.AuditLog{
			UserID:     uid,
			Username:   uname,
			Action:     audit.ActionCreate,
			Resource:   "upload",
			ResourceID: filename,
			Detail:     fmt.Sprintf("上传文件: %s (%d bytes)", header.Filename, header.Size),
			IP:         c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
		})
	}

	response.OK(c, gin.H{"url": url})
}

// isAllowedExt 检查扩展名是否在允许列表中。
func (h *Handler) isAllowedExt(ext string) bool {
	for _, allowed := range h.cfg.AllowedExts {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}

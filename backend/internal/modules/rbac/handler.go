package rbac

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go_web/internal/modules/audit"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/response"
)

// Handler 处理 RBAC 相关 HTTP 请求：用户、角色、权限的 CRUD 及关联操作。
type Handler struct {
	service *Service
	audit   *audit.Service
}

// NewHandler 创建 RBAC Handler。
func NewHandler(service *Service, auditService *audit.Service) *Handler {
	return &Handler{service: service, audit: auditService}
}

// actor 从 Gin Context 提取操作者信息。
func (h *Handler) actor(c *gin.Context) (userID uint, username, ip, userAgent string) {
	if v, ok := c.Get("user_id"); ok {
		userID, _ = v.(uint)
	}
	if v, ok := c.Get("username"); ok {
		username, _ = v.(string)
	}
	return userID, username, c.ClientIP(), c.GetHeader("User-Agent")
}

// logAudit 异步记录操作日志，不受 audit 为 nil 影响。
func (h *Handler) logAudit(userID uint, username, action, resource, resourceID, detail, ip, userAgent string) {
	if h.audit == nil {
		return
	}
	h.audit.Log(&audit.AuditLog{
		UserID:     userID,
		Username:   username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Detail:     detail,
		IP:         ip,
		UserAgent:  userAgent,
	})
}

// ListUsers 分页查询用户列表，支持多条件筛选（用户名/昵称/手机号/邮箱/性别/状态）。
func (h *Handler) ListUsers(c *gin.Context) {
	offset, limit := response.Page(c)
	q := ListUsersQuery{
		Username: c.Query("username"),
		Nickname: c.Query("nickname"),
		Phone:    c.Query("phone"),
		Email:    c.Query("email"),
		Gender:   c.Query("gender"),
		Status:   c.Query("status"),
	}
	users, total, err := h.service.ListUsersPaged(offset, limit, q)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.PaginatedOK(c, users, total, offset, limit)
}

// CreateUser 创建用户，支持同时分配角色。
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if !bind(c, &req) {
		return
	}
	user, err := h.service.CreateUser(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	uid, uname, ip, ua := h.actor(c)
	h.logAudit(uid, uname, audit.ActionCreate, audit.ResourceUser, strconv.FormatUint(uint64(user.ID), 10),
		"创建用户: "+user.Username, ip, ua)
	response.Created(c, user)
}

// GetUser 获取单个用户详情。
func (h *Handler) GetUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	user, err := h.service.GetUser(id)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.OK(c, user)
}

// UpdateUser 更新用户信息（含角色重分配）。
func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdateUserRequest
	if !bind(c, &req) {
		return
	}
	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	uid, uname, ip, ua := h.actor(c)
	h.logAudit(uid, uname, audit.ActionUpdate, audit.ResourceUser, strconv.FormatUint(uint64(id), 10),
		"更新用户: "+user.Username, ip, ua)
	response.OK(c, user)
}

// DeleteUser 删除用户，同时清理 Casbin 分组策略。
func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	uid, uname, ip, ua := h.actor(c)
	if err := h.service.DeleteUser(id); err != nil {
		response.FromError(c, err)
		return
	}
	h.logAudit(uid, uname, audit.ActionDelete, audit.ResourceUser, strconv.FormatUint(uint64(id), 10),
		"删除用户 ID: "+strconv.FormatUint(uint64(id), 10), ip, ua)
	response.OK(c, gin.H{"deleted": true})
}

// AssignRoles 为用户分配角色，覆盖已有角色。
func (h *Handler) AssignRoles(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req AssignRolesRequest
	if !bind(c, &req) {
		return
	}
	uid, uname, ip, ua := h.actor(c)
	if err := h.service.AssignRoles(id, req.RoleIDs); err != nil {
		response.FromError(c, err)
		return
	}
	h.logAudit(uid, uname, audit.ActionUpdate, audit.ResourceUser, strconv.FormatUint(uint64(id), 10),
		"分配角色: "+fmt.Sprintf("%v", req.RoleIDs), ip, ua)
	response.OK(c, gin.H{"assigned": true})
}

// ListRoles 分页查询角色列表，支持按名称/编码/状态筛选。
func (h *Handler) ListRoles(c *gin.Context) {
	offset, limit := response.Page(c)
	q := ListRolesQuery{
		Name:   c.Query("name"),
		Code:   c.Query("code"),
		Status: c.Query("status"),
	}
	roles, total, err := h.service.ListRolesPaged(offset, limit, q)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.PaginatedOK(c, roles, total, offset, limit)
}

// CreateRole 创建角色。
func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if !bind(c, &req) {
		return
	}
	role, err := h.service.CreateRole(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	uid, uname, ip, ua := h.actor(c)
	h.logAudit(uid, uname, audit.ActionCreate, audit.ResourceRole, strconv.FormatUint(uint64(role.ID), 10),
		"创建角色: "+role.Name, ip, ua)
	response.Created(c, role)
}

// UpdateRole 更新角色信息。
func (h *Handler) UpdateRole(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdateRoleRequest
	if !bind(c, &req) {
		return
	}
	role, err := h.service.UpdateRole(id, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	uid, uname, ip, ua := h.actor(c)
	h.logAudit(uid, uname, audit.ActionUpdate, audit.ResourceRole, strconv.FormatUint(uint64(id), 10),
		"更新角色: "+role.Name, ip, ua)
	response.OK(c, role)
}

// DeleteRole 删除角色，同时清理 Casbin 策略。
func (h *Handler) DeleteRole(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	uid, uname, ip, ua := h.actor(c)
	if err := h.service.DeleteRole(id); err != nil {
		response.FromError(c, err)
		return
	}
	h.logAudit(uid, uname, audit.ActionDelete, audit.ResourceRole, strconv.FormatUint(uint64(id), 10),
		"删除角色 ID: "+strconv.FormatUint(uint64(id), 10), ip, ua)
	response.OK(c, gin.H{"deleted": true})
}

// AssignPermissions 为角色分配权限，覆盖已有权限，并同步 Casbin 策略。
func (h *Handler) AssignPermissions(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req AssignPermissionsRequest
	if !bind(c, &req) {
		return
	}
	uid, uname, ip, ua := h.actor(c)
	if err := h.service.AssignPermissions(id, req.PermissionIDs); err != nil {
		response.FromError(c, err)
		return
	}
	h.logAudit(uid, uname, audit.ActionUpdate, audit.ResourceRole, strconv.FormatUint(uint64(id), 10),
		"分配权限: "+fmt.Sprintf("%v", req.PermissionIDs), ip, ua)
	response.OK(c, gin.H{"assigned": true})
}

// ListPermissions 分页查询权限列表，支持按名称/方法筛选。
func (h *Handler) ListPermissions(c *gin.Context) {
	offset, limit := response.Page(c)
	q := ListPermissionsQuery{
		Name:   c.Query("name"),
		Method: c.Query("method"),
	}
	permissions, total, err := h.service.ListPermissionsPaged(offset, limit, q)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.PaginatedOK(c, permissions, total, offset, limit)
}

// CreatePermission 创建权限。
func (h *Handler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if !bind(c, &req) {
		return
	}
	permission, err := h.service.CreatePermission(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	uid, uname, ip, ua := h.actor(c)
	h.logAudit(uid, uname, audit.ActionCreate, audit.ResourcePermission, strconv.FormatUint(uint64(permission.ID), 10),
		"创建权限: "+permission.Name, ip, ua)
	response.Created(c, permission)
}

// UpdatePermission 更新权限，更新后全量同步 Casbin 策略。
func (h *Handler) UpdatePermission(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdatePermissionRequest
	if !bind(c, &req) {
		return
	}
	permission, err := h.service.UpdatePermission(id, req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	uid, uname, ip, ua := h.actor(c)
	h.logAudit(uid, uname, audit.ActionUpdate, audit.ResourcePermission, strconv.FormatUint(uint64(id), 10),
		"更新权限: "+permission.Name, ip, ua)
	response.OK(c, permission)
}

// DeletePermission 删除权限，删除后全量同步 Casbin 策略。
func (h *Handler) DeletePermission(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	uid, uname, ip, ua := h.actor(c)
	if err := h.service.DeletePermission(id); err != nil {
		response.FromError(c, err)
		return
	}
	h.logAudit(uid, uname, audit.ActionDelete, audit.ResourcePermission, strconv.FormatUint(uint64(id), 10),
		"删除权限 ID: "+strconv.FormatUint(uint64(id), 10), ip, ua)
	response.OK(c, gin.H{"deleted": true})
}

// BatchDeleteUsers 批量删除用户。
func (h *Handler) BatchDeleteUsers(c *gin.Context) {
	var req BatchDeleteRequest
	if !bind(c, &req) {
		return
	}
	uid, uname, ip, ua := h.actor(c)
	if err := h.service.BatchDeleteUsers(req.IDs); err != nil {
		response.FromError(c, err)
		return
	}
	h.logAudit(uid, uname, audit.ActionDelete, audit.ResourceUser, "batch",
		"批量删除用户 IDs: "+fmt.Sprintf("%v", req.IDs), ip, ua)
	response.OK(c, gin.H{"deleted": true})
}

// BatchDeleteRoles 批量删除角色。
func (h *Handler) BatchDeleteRoles(c *gin.Context) {
	var req BatchDeleteRequest
	if !bind(c, &req) {
		return
	}
	uid, uname, ip, ua := h.actor(c)
	if err := h.service.BatchDeleteRoles(req.IDs); err != nil {
		response.FromError(c, err)
		return
	}
	h.logAudit(uid, uname, audit.ActionDelete, audit.ResourceRole, "batch",
		"批量删除角色 IDs: "+fmt.Sprintf("%v", req.IDs), ip, ua)
	response.OK(c, gin.H{"deleted": true})
}

// BatchDeletePermissions 批量删除权限。
func (h *Handler) BatchDeletePermissions(c *gin.Context) {
	var req BatchDeleteRequest
	if !bind(c, &req) {
		return
	}
	uid, uname, ip, ua := h.actor(c)
	if err := h.service.BatchDeletePermissions(req.IDs); err != nil {
		response.FromError(c, err)
		return
	}
	h.logAudit(uid, uname, audit.ActionDelete, audit.ResourcePermission, "batch",
		"批量删除权限 IDs: "+fmt.Sprintf("%v", req.IDs), ip, ua)
	response.OK(c, gin.H{"deleted": true})
}

// bind 解析 JSON 请求体并校验，校验失败返回 400 错误响应。
func bind(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.Error(c, 400, apperror.CodeInvalidArgument, err.Error())
		return false
	}
	return true
}

// parseID 从路径参数中解析 uint 类型的 ID，解析失败返回 400。
func parseID(c *gin.Context) (uint, bool) {
	value := c.Param("id")
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, 400, apperror.CodeInvalidArgument, "invalid id")
		return 0, false
	}
	return uint(id), true
}

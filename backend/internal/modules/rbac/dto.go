package rbac

// CreateUserRequest 创建用户请求体。
type CreateUserRequest struct {
	Username string     `json:"username" binding:"required,min=2,max=64"`
	Nickname string     `json:"nickname"`
	Email    string     `json:"email" binding:"required,email,max=128"`
	Phone    string     `json:"phone"`
	Gender   string     `json:"gender"`
	Password string     `json:"password" binding:"required,min=6,max=128"`
	Status   UserStatus `json:"status"`
	RoleIDs  []uint     `json:"role_ids"` // 初始角色 ID 列表
}

// UpdateUserRequest 更新用户请求体。
type UpdateUserRequest struct {
	Username string     `json:"username" binding:"required,min=2,max=64"`
	Nickname string     `json:"nickname"`
	Email    string     `json:"email" binding:"required,email,max=128"`
	Phone    string     `json:"phone"`
	Gender   string     `json:"gender"`
	Status   UserStatus `json:"status" binding:"required"`
	RoleIDs  []uint     `json:"role_ids"`
}

// CreateRoleRequest 创建角色请求体。
type CreateRoleRequest struct {
	Code        string     `json:"code" binding:"required,min=2,max=64"`
	Name        string     `json:"name" binding:"required,min=2,max=128"`
	Description string     `json:"description" binding:"max=255"`
	Status      UserStatus `json:"status"`
}

// UpdateRoleRequest 更新角色请求体。
type UpdateRoleRequest struct {
	Name        string     `json:"name" binding:"required,min=2,max=128"`
	Description string     `json:"description" binding:"max=255"`
	Status      UserStatus `json:"status"`
}

// CreatePermissionRequest 创建权限请求体。
type CreatePermissionRequest struct {
	Code        string `json:"code" binding:"required,min=2,max=128"`
	Name        string `json:"name" binding:"required,min=2,max=128"`
	Method      string `json:"method" binding:"required,max=16"`   // HTTP 方法（GET/POST/PUT/DELETE）
	Path        string `json:"path" binding:"required,max=255"`    // 路由路径模板
	Description string `json:"description" binding:"max=255"`
}

// UpdatePermissionRequest 更新权限请求体。
type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=128"`
	Method      string `json:"method" binding:"required,max=16"`
	Path        string `json:"path" binding:"required,max=255"`
	Description string `json:"description" binding:"max=255"`
}

// AssignRolesRequest 为用户分配角色请求体。
type AssignRolesRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// AssignPermissionsRequest 为角色分配权限请求体。
type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// BatchDeleteRequest 批量删除请求体。
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

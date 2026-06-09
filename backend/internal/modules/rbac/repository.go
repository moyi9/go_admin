package rbac

import (
	"errors"

	"gorm.io/gorm"
)

// Repository 封装 RBAC 模块的数据库操作。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建 Repository。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

// ListUsersQuery 用户列表多条件查询参数。
type ListUsersQuery struct {
	Username string
	Nickname string
	Phone    string
	Email    string
	Gender   string
	Status   string
	Keyword  string // 原始 keyword 搜索（用户名/邮箱）
}

// ListUsers 查询所有用户（含角色），按 ID 升序排列。
func (r *Repository) ListUsers() ([]User, error) {
	var users []User
	err := r.db.Preload("Roles").Order("id asc").Find(&users).Error
	return users, err
}

// CountUsers 返回用户总数。
func (r *Repository) CountUsers(q ListUsersQuery) (int64, error) {
	var count int64
	query := r.db.Model(&User{})
	query = applyUserFilters(query, q)
	return count, query.Count(&count).Error
}

// ListUsersPaged 分页查询用户。
func (r *Repository) ListUsersPaged(offset, limit int, q ListUsersQuery) ([]User, error) {
	var users []User
	query := r.db.Preload("Roles").Order("id asc")
	query = applyUserFilters(query, q)
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// applyUserFilters 根据 ListUsersQuery 动态构建 WHERE 条件。
func applyUserFilters(db *gorm.DB, q ListUsersQuery) *gorm.DB {
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		db = db.Where("username LIKE ? OR email LIKE ?", like, like)
	}
	if q.Username != "" {
		db = db.Where("username LIKE ?", "%"+q.Username+"%")
	}
	if q.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+q.Nickname+"%")
	}
	if q.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+q.Phone+"%")
	}
	if q.Email != "" {
		db = db.Where("email LIKE ?", "%"+q.Email+"%")
	}
	if q.Gender != "" {
		db = db.Where("gender = ?", q.Gender)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	return db
}

// GetUser 根据 ID 查询用户（含角色和权限）。
func (r *Repository) GetUser(id uint) (*User, error) {
	var user User
	err := r.db.Preload("Roles.Permissions").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetUserByUsername 根据用户名查询用户（含角色和权限），用于登录校验。
func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// CreateUser 创建用户。
func (r *Repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

// UpdateUser 更新用户字段，使用 map 仅更新非零值字段。
func (r *Repository) UpdateUser(user *User) error {
	return r.db.Model(user).Updates(map[string]interface{}{
		"username":   user.Username,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"phone":      user.Phone,
		"gender":     user.Gender,
		"avatar_url": user.AvatarURL,
		"status":     user.Status,
	}).Error
}

// UpdateUserPassword 更新用户密码哈希。
func (r *Repository) UpdateUserPassword(id uint, passwordHash string) error {
	return r.db.Model(&User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
}

// DeleteUser 软删除用户。
func (r *Repository) DeleteUser(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

// DeleteUsers 批量软删除用户。
func (r *Repository) DeleteUsers(ids []uint) error {
	return r.db.Delete(&User{}, "id IN ?", ids).Error
}

// SetUserRoles 替换用户的角色关联（全量覆盖）。
func (r *Repository) SetUserRoles(user *User, roles []Role) error {
	return r.db.Model(user).Association("Roles").Replace(roles)
}

// ListRoles 查询所有角色（含权限），按 ID 升序排列。
func (r *Repository) ListRoles() ([]Role, error) {
	var roles []Role
	err := r.db.Preload("Permissions").Order("id asc").Find(&roles).Error
	return roles, err
}

// ListRolesQuery 角色列表多条件查询参数。
type ListRolesQuery struct {
	Name   string
	Code   string
	Status string
}

// CountRoles 返回角色总数。
func (r *Repository) CountRoles(q ListRolesQuery) (int64, error) {
	var count int64
	query := r.db.Model(&Role{})
	query = applyRoleFilters(query, q)
	return count, query.Count(&count).Error
}

// ListRolesPaged 分页查询角色。
func (r *Repository) ListRolesPaged(offset, limit int, q ListRolesQuery) ([]Role, error) {
	var roles []Role
	query := r.db.Preload("Permissions").Order("id asc")
	query = applyRoleFilters(query, q)
	err := query.Offset(offset).Limit(limit).Find(&roles).Error
	return roles, err
}

// applyRoleFilters 根据 ListRolesQuery 动态构建 WHERE 条件。
func applyRoleFilters(db *gorm.DB, q ListRolesQuery) *gorm.DB {
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.Code != "" {
		db = db.Where("code LIKE ?", "%"+q.Code+"%")
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	return db
}

// GetRole 根据 ID 查询角色（含权限）。
func (r *Repository) GetRole(id uint) (*Role, error) {
	var role Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

// GetRoleByCode 根据编码查询角色（含权限）。
func (r *Repository) GetRoleByCode(code string) (*Role, error) {
	var role Role
	err := r.db.Preload("Permissions").Where("code = ?", code).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

// GetRolesByIDs 根据 ID 列表批量查询角色。
func (r *Repository) GetRolesByIDs(ids []uint) ([]Role, error) {
	var roles []Role
	err := r.db.Where("id IN ?", ids).Find(&roles).Error
	return roles, err
}

// CreateRole 创建角色。
func (r *Repository) CreateRole(role *Role) error {
	return r.db.Create(role).Error
}

// UpdateRole 更新角色的名称和描述。
func (r *Repository) UpdateRole(role *Role) error {
	return r.db.Model(role).Updates(map[string]interface{}{"name": role.Name, "description": role.Description}).Error
}

// DeleteRole 软删除角色。
func (r *Repository) DeleteRole(id uint) error {
	return r.db.Delete(&Role{}, id).Error
}

// DeleteRoles 批量删除角色。
func (r *Repository) DeleteRoles(ids []uint) error {
	return r.db.Delete(&Role{}, "id IN ?", ids).Error
}

// SetRolePermissions 替换角色的权限关联（全量覆盖）。
func (r *Repository) SetRolePermissions(role *Role, permissions []Permission) error {
	return r.db.Model(role).Association("Permissions").Replace(permissions)
}

// ListPermissions 查询所有权限，按 ID 升序排列。
func (r *Repository) ListPermissions() ([]Permission, error) {
	var permissions []Permission
	err := r.db.Order("id asc").Find(&permissions).Error
	return permissions, err
}

// ListPermissionsQuery 权限列表多条件查询参数。
type ListPermissionsQuery struct {
	Name   string
	Method string
}

// CountPermissions 返回权限总数。
func (r *Repository) CountPermissions(q ListPermissionsQuery) (int64, error) {
	var count int64
	query := r.db.Model(&Permission{})
	query = applyPermissionFilters(query, q)
	return count, query.Count(&count).Error
}

// ListPermissionsPaged 分页查询权限。
func (r *Repository) ListPermissionsPaged(offset, limit int, q ListPermissionsQuery) ([]Permission, error) {
	var permissions []Permission
	query := r.db.Order("id asc")
	query = applyPermissionFilters(query, q)
	err := query.Offset(offset).Limit(limit).Find(&permissions).Error
	return permissions, err
}

// applyPermissionFilters 根据 ListPermissionsQuery 动态构建 WHERE 条件。
func applyPermissionFilters(db *gorm.DB, q ListPermissionsQuery) *gorm.DB {
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.Method != "" {
		db = db.Where("method = ?", q.Method)
	}
	return db
}

// GetPermission 根据 ID 查询权限。
func (r *Repository) GetPermission(id uint) (*Permission, error) {
	var permission Permission
	err := r.db.First(&permission, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &permission, err
}

// GetPermissionsByIDs 根据 ID 列表批量查询权限。
func (r *Repository) GetPermissionsByIDs(ids []uint) ([]Permission, error) {
	var permissions []Permission
	err := r.db.Where("id IN ?", ids).Find(&permissions).Error
	return permissions, err
}

// GetPermissionByCode 根据编码查询权限。
func (r *Repository) GetPermissionByCode(code string) (*Permission, error) {
	var permission Permission
	err := r.db.Where("code = ?", code).First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &permission, err
}

// CreatePermission 创建权限。
func (r *Repository) CreatePermission(permission *Permission) error {
	return r.db.Create(permission).Error
}

// UpdatePermission 更新权限字段。
func (r *Repository) UpdatePermission(permission *Permission) error {
	return r.db.Model(permission).Updates(map[string]interface{}{"name": permission.Name, "method": permission.Method, "path": permission.Path, "description": permission.Description}).Error
}

// DeletePermission 软删除权限。
func (r *Repository) DeletePermission(id uint) error {
	return r.db.Delete(&Permission{}, id).Error
}

// DeletePermissions 批量删除权限。
func (r *Repository) DeletePermissions(ids []uint) error {
	return r.db.Delete(&Permission{}, "id IN ?", ids).Error
}

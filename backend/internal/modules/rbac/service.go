package rbac

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/security"
)

type Service struct {
	repo     *Repository
	enforcer *casbin.Enforcer
}

func NewService(repo *Repository, enforcer *casbin.Enforcer) *Service {
	return &Service{repo: repo, enforcer: enforcer}
}

// ListUsers 查询全部用户（无分页），用于全量策略同步。
func (s *Service) ListUsers() ([]User, error) {
	return s.repo.ListUsers()
}

// ListUsersPaged 分页查询用户，返回数据、总数。
func (s *Service) ListUsersPaged(offset, limit int, q ListUsersQuery) ([]User, int64, error) {
	total, err := s.repo.CountUsers(q)
	if err != nil {
		return nil, 0, err
	}
	users, err := s.repo.ListUsersPaged(offset, limit, q)
	return users, total, err
}

// GetUser 根据 ID 获取用户，不存在返回 404 错误。
func (s *Service) GetUser(id uint) (*User, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.New(apperror.CodeNotFound, "user not found")
	}
	return user, nil
}

func (s *Service) CreateUser(req CreateUserRequest) (*User, error) {
	status := req.Status
	if status == "" {
		status = UserStatusActive
	}
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username:     req.Username,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Phone:        req.Phone,
		Gender:       req.Gender,
		PasswordHash: hash,
		Status:       status,
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, apperror.Wrap(apperror.CodeConflict, "user already exists or invalid", err)
	}
	if len(req.RoleIDs) > 0 {
		roles, err := s.repo.GetRolesByIDs(req.RoleIDs)
		if err != nil {
			return nil, err
		}
		if err := s.repo.SetUserRoles(user, roles); err != nil {
			return nil, err
		}
		if err := s.syncUserPolicies(user.ID); err != nil {
			return nil, err
		}
	}
	return user, nil
}

// UpdateUser 更新用户信息。如果传入了 role_ids，会重新分配用户角色并同步 Casbin 策略。
func (s *Service) UpdateUser(id uint, req UpdateUserRequest) (*User, error) {
	user, err := s.GetUser(id)
	if err != nil {
		return nil, err
	}
	user.Username = req.Username
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Gender = req.Gender
	user.Status = req.Status
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, apperror.Wrap(apperror.CodeConflict, "user already exists or invalid", err)
	}
	if len(req.RoleIDs) > 0 {
		roles, err := s.repo.GetRolesByIDs(req.RoleIDs)
		if err != nil {
			return nil, err
		}
		if err := s.repo.SetUserRoles(user, roles); err != nil {
			return nil, err
		}
		if err := s.syncUserPolicies(user.ID); err != nil {
			return nil, err
		}
	}
	return user, nil
}

// DeleteUser 删除用户并从 Casbin 中移除分组策略。
func (s *Service) DeleteUser(id uint) error {
	if _, err := s.GetUser(id); err != nil {
		return err
	}
	if err := s.repo.DeleteUser(id); err != nil {
		return err
	}
	if s.enforcer != nil {
		s.enforcer.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("%d", id))
	}
	return nil
}

// BatchDeleteUsers 批量删除用户。
func (s *Service) BatchDeleteUsers(ids []uint) error {
	if err := s.repo.DeleteUsers(ids); err != nil {
		return err
	}
	if s.enforcer != nil {
		for _, id := range ids {
			s.enforcer.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("%d", id))
		}
	}
	return nil
}

// AssignRoles 为用户分配角色（覆盖），并同步 Casbin 分组策略。
func (s *Service) AssignRoles(userID uint, roleIDs []uint) error {
	user, err := s.GetUser(userID)
	if err != nil {
		return err
	}
	roles, err := s.repo.GetRolesByIDs(roleIDs)
	if err != nil {
		return err
	}
	if len(roles) != len(roleIDs) {
		return apperror.New(apperror.CodeInvalidArgument, "some roles do not exist")
	}
	if err := s.repo.SetUserRoles(user, roles); err != nil {
		return err
	}
	return s.syncUserPolicies(userID)
}

// ListRoles 查询全部角色（无分页），用于全量策略同步。
func (s *Service) ListRoles() ([]Role, error) {
	return s.repo.ListRoles()
}

// ListRolesPaged 分页查询角色，返回数据、总数。
func (s *Service) ListRolesPaged(offset, limit int, q ListRolesQuery) ([]Role, int64, error) {
	total, err := s.repo.CountRoles(q)
	if err != nil {
		return nil, 0, err
	}
	roles, err := s.repo.ListRolesPaged(offset, limit, q)
	return roles, total, err
}

// CreateRole 创建角色，状态为空时默认启用。
func (s *Service) CreateRole(req CreateRoleRequest) (*Role, error) {
	status := req.Status
	if status == "" {
		status = UserStatusActive
	}
	role := &Role{Code: req.Code, Name: req.Name, Description: req.Description, Status: status}
	if err := s.repo.CreateRole(role); err != nil {
		return nil, apperror.Wrap(apperror.CodeConflict, "role already exists or invalid", err)
	}
	return role, nil
}

// UpdateRole 更新角色信息。状态字段非空时才更新，避免覆盖已有值。
func (s *Service) UpdateRole(id uint, req UpdateRoleRequest) (*Role, error) {
	role, err := s.repo.GetRole(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, apperror.New(apperror.CodeNotFound, "role not found")
	}
	role.Name = req.Name
	role.Description = req.Description
	if req.Status != "" {
		role.Status = req.Status
	}
	if err := s.repo.UpdateRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

// DeleteRole 删除角色并清理 Casbin 中对应的策略和分组。
func (s *Service) DeleteRole(id uint) error {
	role, err := s.repo.GetRole(id)
	if err != nil {
		return err
	}
	if role == nil {
		return apperror.New(apperror.CodeNotFound, "role not found")
	}
	if err := s.repo.DeleteRole(id); err != nil {
		return err
	}
	if s.enforcer != nil {
		s.enforcer.RemoveFilteredPolicy(0, role.Code)
		s.enforcer.RemoveFilteredGroupingPolicy(1, role.Code)
	}
	return nil
}

// BatchDeleteRoles 批量删除角色。
func (s *Service) BatchDeleteRoles(ids []uint) error {
	roles, err := s.repo.GetRolesByIDs(ids)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteRoles(ids); err != nil {
		return err
	}
	if s.enforcer != nil {
		for _, role := range roles {
			s.enforcer.RemoveFilteredPolicy(0, role.Code)
			s.enforcer.RemoveFilteredGroupingPolicy(1, role.Code)
		}
	}
	return nil
}

// AssignPermissions 为角色分配权限（覆盖），并同步 Casbin 策略。
func (s *Service) AssignPermissions(roleID uint, permissionIDs []uint) error {
	role, err := s.repo.GetRole(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return apperror.New(apperror.CodeNotFound, "role not found")
	}
	permissions, err := s.repo.GetPermissionsByIDs(permissionIDs)
	if err != nil {
		return err
	}
	if len(permissions) != len(permissionIDs) {
		return apperror.New(apperror.CodeInvalidArgument, "some permissions do not exist")
	}
	if err := s.repo.SetRolePermissions(role, permissions); err != nil {
		return err
	}
	return s.syncRolePolicies(role.Code)
}

// ListPermissions 查询全部权限（无分页），用于全量策略同步。
func (s *Service) ListPermissions() ([]Permission, error) {
	return s.repo.ListPermissions()
}

// ListPermissionsPaged 分页查询权限，返回数据、总数。
func (s *Service) ListPermissionsPaged(offset, limit int, q ListPermissionsQuery) ([]Permission, int64, error) {
	total, err := s.repo.CountPermissions(q)
	if err != nil {
		return nil, 0, err
	}
	perms, err := s.repo.ListPermissionsPaged(offset, limit, q)
	return perms, total, err
}

// CreatePermission 创建权限。
func (s *Service) CreatePermission(req CreatePermissionRequest) (*Permission, error) {
	permission := &Permission{
		Code:        req.Code,
		Name:        req.Name,
		Method:      strings.ToUpper(req.Method),
		Path:        req.Path,
		Description: req.Description,
	}
	if err := s.repo.CreatePermission(permission); err != nil {
		return nil, apperror.Wrap(apperror.CodeConflict, "permission already exists or invalid", err)
	}
	return permission, nil
}

// UpdatePermission 更新权限，更新后全量同步 Casbin 策略。
func (s *Service) UpdatePermission(id uint, req UpdatePermissionRequest) (*Permission, error) {
	permission, err := s.repo.GetPermission(id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, apperror.New(apperror.CodeNotFound, "permission not found")
	}
	permission.Name = req.Name
	permission.Method = strings.ToUpper(req.Method)
	permission.Path = req.Path
	permission.Description = req.Description
	if err := s.repo.UpdatePermission(permission); err != nil {
		return nil, err
	}
	if err := s.SyncPolicies(); err != nil {
		return nil, err
	}
	return permission, nil
}

// DeletePermission 删除权限，删除后全量同步 Casbin 策略。
func (s *Service) DeletePermission(id uint) error {
	permission, err := s.repo.GetPermission(id)
	if err != nil {
		return err
	}
	if permission == nil {
		return apperror.New(apperror.CodeNotFound, "permission not found")
	}
	if err := s.repo.DeletePermission(id); err != nil {
		return err
	}
	return s.SyncPolicies()
}

// BatchDeletePermissions 批量删除权限。
func (s *Service) BatchDeletePermissions(ids []uint) error {
	if err := s.repo.DeletePermissions(ids); err != nil {
		return err
	}
	return s.SyncPolicies()
}

// syncUserPolicies 删除该用户的旧 Casbin 分组策略，并根据当前角色重新添加。
func (s *Service) syncUserPolicies(userID uint) error {
	if s.enforcer == nil {
		return nil
	}
	userIDStr := fmt.Sprintf("%d", userID)
	s.enforcer.RemoveFilteredGroupingPolicy(0, userIDStr)
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	for _, role := range user.Roles {
		if _, err := s.enforcer.AddGroupingPolicy(userIDStr, role.Code); err != nil {
			return fmt.Errorf("add user role policy: %w", err)
		}
	}
	return nil
}

// syncRolePolicies 删除该角色的旧 Casbin 策略，并根据当前权限重新添加。
func (s *Service) syncRolePolicies(roleCode string) error {
	if s.enforcer == nil {
		return nil
	}
	s.enforcer.RemoveFilteredPolicy(0, roleCode)
	role, err := s.repo.GetRoleByCode(roleCode)
	if err != nil {
		return err
	}
	if role == nil {
		return nil
	}
	for _, perm := range role.Permissions {
		if _, err := s.enforcer.AddPolicy(roleCode, perm.Path, strings.ToUpper(perm.Method)); err != nil {
			return fmt.Errorf("add role policy: %w", err)
		}
	}
	return nil
}

// SyncPolicies 全量同步 Casbin 策略：清除旧策略后，根据数据库中所有角色和用户的关联关系重新构建。
func (s *Service) SyncPolicies() error {
	if s.enforcer == nil {
		return nil
	}
	s.enforcer.ClearPolicy()

	roles, err := s.repo.ListRoles()
	if err != nil {
		return err
	}
	for _, role := range roles {
		for _, permission := range role.Permissions {
			if _, err := s.enforcer.AddPolicy(role.Code, permission.Path, strings.ToUpper(permission.Method)); err != nil {
				return fmt.Errorf("add role policy: %w", err)
			}
		}
	}

	users, err := s.repo.ListUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		for _, role := range user.Roles {
			if _, err := s.enforcer.AddGroupingPolicy(fmt.Sprintf("%d", user.ID), role.Code); err != nil {
				return fmt.Errorf("add user role policy: %w", err)
			}
		}
	}

	return nil
}

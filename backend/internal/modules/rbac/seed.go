package rbac

import (
	"strings"

	"go_web/internal/config"
	"go_web/internal/pkg/security"
	"gorm.io/gorm"
)

// Seed 初始化开发环境默认管理员、管理员角色和接口权限。
// 生产环境建议关闭 seed.enabled，或改为一次性的显式初始化流程。
func Seed(db *gorm.DB, cfg config.SeedConfig) error {
	if !cfg.Enabled {
		return nil
	}

	adminRole := Role{Code: "admin", Name: "Administrator", Description: "Built-in administrator role", Status: UserStatusActive}
	if err := db.FirstOrCreate(&adminRole, Role{Code: adminRole.Code}).Error; err != nil {
		return err
	}

	permissions := defaultPermissions()
	// 先保存中文名称映射，因为 FirstOrCreate 会将 struct 覆写为数据库中的旧值
	chineseNames := make(map[string]string, len(permissions))
	for _, p := range permissions {
		chineseNames[p.Code] = p.Name
	}
	for i := range permissions {
		if err := db.FirstOrCreate(&permissions[i], Permission{Code: permissions[i].Code}).Error; err != nil {
			return err
		}
	}
	// 用预存的中文名称更新已有记录（FirstOrCreate 不会更新非零值字段）
	for code, cnName := range chineseNames {
		db.Model(&Permission{}).Where("code = ?", code).Update("name", cnName)
	}
	if err := db.Model(&adminRole).Association("Permissions").Replace(permissions); err != nil {
		return err
	}

	passwordHash, err := security.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}
	admin := User{
		Username:     cfg.AdminUsername,
		Email:        cfg.AdminEmail,
		PasswordHash: passwordHash,
		Status:       UserStatusActive,
	}
	if err := db.Where(User{Username: admin.Username}).Attrs(admin).FirstOrCreate(&admin).Error; err != nil {
		return err
	}
	return db.Model(&admin).Association("Roles").Replace([]Role{adminRole})
}

// defaultPermissions 覆盖脚手架内置的受保护接口。
// 新增业务模块后，可以按相同格式追加默认权限，或通过权限 API 动态创建。
func defaultPermissions() []Permission {
	routes := []struct {
		method string
		path   string
		name   string
	}{
		{"GET", "/api/v1/auth/me", "查看当前用户"},
		{"GET", "/api/v1/users", "用户列表"},
		{"POST", "/api/v1/users", "创建用户"},
		{"GET", "/api/v1/users/:id", "查看用户详情"},
		{"PUT", "/api/v1/users/:id", "更新用户"},
		{"DELETE", "/api/v1/users/:id", "删除用户"},
		{"DELETE", "/api/v1/users/batch", "批量删除用户"},
		{"POST", "/api/v1/users/:id/roles", "分配用户角色"},
		{"GET", "/api/v1/roles", "角色列表"},
		{"POST", "/api/v1/roles", "创建角色"},
		{"PUT", "/api/v1/roles/:id", "更新角色"},
		{"DELETE", "/api/v1/roles/:id", "删除角色"},
		{"DELETE", "/api/v1/roles/batch", "批量删除角色"},
		{"POST", "/api/v1/roles/:id/permissions", "分配角色权限"},
		{"GET", "/api/v1/permissions", "权限列表"},
		{"POST", "/api/v1/permissions", "创建权限"},
		{"PUT", "/api/v1/permissions/:id", "更新权限"},
		{"DELETE", "/api/v1/permissions/:id", "删除权限"},
		{"DELETE", "/api/v1/permissions/batch", "批量删除权限"},
		{"POST", "/api/v1/upload", "文件上传"},
		{"GET", "/api/v1/audit-logs", "操作日志"},
		{"GET", "/api/v1/notifications", "通知列表"},
		{"GET", "/api/v1/notifications/unread", "未读通知数"},
		{"GET", "/api/v1/notifications/unread-list", "未读通知列表"},
		{"POST", "/api/v1/notifications", "发送通知"},
		{"PUT", "/api/v1/notifications/:id/read", "标记已读"},
		{"PUT", "/api/v1/notifications/read-all", "全部已读"},
		{"GET", "/api/v1/dashboard/stats", "仪表盘统计"},
	}
	permissions := make([]Permission, 0, len(routes))
	for _, route := range routes {
		code := strings.ToLower(route.method) + "." + strings.TrimLeft(strings.ReplaceAll(route.path, "/", "."), ".")
		permissions = append(permissions, Permission{
			Code:   code,
			Name:   route.name,
			Method: route.method,
			Path:   route.path,
		})
	}
	return permissions
}

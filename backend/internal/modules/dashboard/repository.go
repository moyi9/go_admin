package dashboard

import (
	"go_web/internal/modules/rbac"
	"gorm.io/gorm"
)

// Repository 封装仪表盘统计的数据库查询。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GenderDistribution 按性别分组统计用户数。
func (r *Repository) GenderDistribution() (map[string]int64, error) {
	type result struct {
		Gender string
		Count  int64
	}
	var rows []result
	if err := r.db.Model(&rbac.User{}).
		Select("gender, count(*) as count").
		Group("gender").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, row := range rows {
		out[row.Gender] = row.Count
	}
	return out, nil
}

// StatusDistribution 按状态分组统计用户数。
func (r *Repository) StatusDistribution() (map[string]int64, error) {
	type result struct {
		Status string
		Count  int64
	}
	var rows []result
	if err := r.db.Model(&rbac.User{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, row := range rows {
		out[row.Status] = row.Count
	}
	return out, nil
}

// RoleDistribution 统计每个角色关联的用户数。
func (r *Repository) RoleDistribution() ([]RoleCount, error) {
	type result struct {
		RoleName string
		Count    int64
	}
	var rows []result
	if err := r.db.Table("roles").
		Select("roles.name as role_name, count(user_roles.user_id) as count").
		Joins("left join user_roles on user_roles.role_id = roles.id").
		Group("roles.id, roles.name").
		Order("count desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]RoleCount, len(rows))
	for i, row := range rows {
		out[i] = RoleCount{RoleName: row.RoleName, UserCount: row.Count}
	}
	return out, nil
}

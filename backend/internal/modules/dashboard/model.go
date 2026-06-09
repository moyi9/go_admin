package dashboard

// Stats 仪表盘统计聚合响应。
type Stats struct {
	GenderDistribution map[string]int64 `json:"gender_distribution"` // male/female/unknown → 人数
	StatusDistribution map[string]int64 `json:"status_distribution"` // active/disabled → 人数
	RoleDistribution   []RoleCount      `json:"role_distribution"`   // 每个角色的用户数
}

// RoleCount 角色-用户数统计项。
type RoleCount struct {
	RoleName   string `json:"role_name"`
	UserCount  int64  `json:"user_count"`
}

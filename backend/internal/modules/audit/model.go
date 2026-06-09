// Package audit 提供操作日志记录与查询功能。
package audit

import "time"

// ActionType 操作类型常量。
const (
	ActionLogin  = "LOGIN"
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)

// ResourceType 资源类型常量。
const (
	ResourceUser       = "user"
	ResourceRole       = "role"
	ResourcePermission = "permission"
)

// AuditLog 操作日志，记录用户在系统中的关键操作。
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Username   string    `gorm:"size:64;not null" json:"username"`
	Action     string    `gorm:"size:32;not null;index" json:"action"`
	Resource   string    `gorm:"size:32;not null;index" json:"resource"`
	ResourceID string    `gorm:"size:32" json:"resource_id"`
	Detail     string    `gorm:"size:512" json:"detail"`
	IP         string    `gorm:"size:64" json:"ip"`
	UserAgent  string    `gorm:"size:512" json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 自定义表名。
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Models 返回需要 AutoMigrate 的模型列表。
func Models() []any {
	return []any{&AuditLog{}}
}

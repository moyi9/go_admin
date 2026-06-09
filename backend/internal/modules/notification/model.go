// Package notification 提供站内通知功能，支持系统自动通知和管理员手动发送。
package notification

import "time"

// NotificationType 通知类型常量。
const (
	TypeSystem   = "system"   // 系统通知（资料更新等）
	TypeSecurity = "security" // 安全提醒（密码修改等）
)

// Notification 站内通知，定向发送给指定用户。
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Type      string    `gorm:"size:32;not null" json:"type"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"size:2000" json:"content"`
	Link      string    `gorm:"size:512" json:"link"`
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 自定义表名。
func (Notification) TableName() string {
	return "notifications"
}

// Models 返回需要 AutoMigrate 的模型列表。
func Models() []any {
	return []any{&Notification{}}
}

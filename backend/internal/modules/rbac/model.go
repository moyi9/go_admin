package rbac

import "time"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

// User 是后台账号。用户通过 user_roles 关联多个角色。
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Nickname     string     `gorm:"size:64" json:"nickname"`
	Email        string     `gorm:"size:128;uniqueIndex;not null" json:"email"`
	Phone        string     `gorm:"size:20" json:"phone"`
	Gender       string     `gorm:"size:8;default:unknown" json:"gender"`
	AvatarURL    string     `gorm:"size:512" json:"avatar_url"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Status       UserStatus `gorm:"size:32;not null;default:active" json:"status"`
	Roles        []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Role 是权限集合。角色通过 role_permissions 关联多个接口权限。
type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Code        string       `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name        string       `gorm:"size:128;not null" json:"name"`
	Description string       `gorm:"size:255" json:"description"`
	Status      UserStatus   `gorm:"size:32;not null;default:active" json:"status"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Permission 表示一个可授权的接口权限，使用 Method + Path 作为 Casbin 策略对象。
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:128;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Method      string    `gorm:"size:16;not null" json:"method"`
	Path        string    `gorm:"size:255;not null" json:"path"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRole 是用户和角色的多对多关系表。
type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

// RolePermission 是角色和权限的多对多关系表。
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

// Models 集中返回需要 AutoMigrate 的 RBAC 模型。
func Models() []any {
	return []any{
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
	}
}

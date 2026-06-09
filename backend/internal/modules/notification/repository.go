package notification

import (
	"gorm.io/gorm"
)

// Repository 封装通知的数据库操作。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建 Repository。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 写入一条通知。
func (r *Repository) Create(n *Notification) error {
	return r.db.Create(n).Error
}

// ListByUser 分页查询指定用户的通知，按时间倒序排列。
func (r *Repository) ListByUser(userID uint, offset, limit int) ([]Notification, int64, error) {
	var list []Notification
	var total int64
	query := r.db.Model(&Notification{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&list).Error
	if list == nil {
		list = []Notification{}
	}
	return list, total, err
}

// ListUnread 查询指定用户的最新 N 条未读通知。
func (r *Repository) ListUnread(userID uint, limit int) ([]Notification, error) {
	var list []Notification
	err := r.db.Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at desc").
		Limit(limit).
		Find(&list).Error
	if list == nil {
		list = []Notification{}
	}
	return list, err
}

// CountUnread 返回指定用户的未读通知数。
func (r *Repository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead 将指定用户的一条通知标记为已读。
func (r *Repository) MarkAsRead(userID, id uint) error {
	return r.db.Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead 将指定用户的全部通知标记为已读。
func (r *Repository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

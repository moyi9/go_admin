package audit

import (
	"gorm.io/gorm"
)

// Repository 封装操作日志的数据库操作。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建 Repository。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Insert 写入一条操作日志。
func (r *Repository) Insert(log *AuditLog) error {
	return r.db.Create(log).Error
}

// ListAuditLogsQuery 日志列表查询参数。
type ListAuditLogsQuery struct {
	Action   string
	Resource string
	Keyword  string // 搜索用户名或 detail
}

// Count 返回符合条件的日志总数。
func (r *Repository) Count(q ListAuditLogsQuery, since, until string) (int64, error) {
	var count int64
	query := r.db.Model(&AuditLog{})
	query = applyFilters(query, q, since, until)
	return count, query.Count(&count).Error
}

// ListPaged 分页查询操作日志，按时间倒序排列。
func (r *Repository) ListPaged(offset, limit int, q ListAuditLogsQuery, since, until string) ([]AuditLog, error) {
	var logs []AuditLog
	query := r.db.Model(&AuditLog{}).Order("created_at desc")
	query = applyFilters(query, q, since, until)
	err := query.Offset(offset).Limit(limit).Find(&logs).Error
	return logs, err
}

// applyFilters 根据查询条件动态构建 WHERE 子句。
func applyFilters(db *gorm.DB, q ListAuditLogsQuery, since, until string) *gorm.DB {
	if q.Action != "" {
		db = db.Where("action = ?", q.Action)
	}
	if q.Resource != "" {
		db = db.Where("resource = ?", q.Resource)
	}
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		db = db.Where("username LIKE ? OR detail LIKE ?", like, like)
	}
	if since != "" {
		db = db.Where("created_at >= ?", since)
	}
	if until != "" {
		db = db.Where("created_at <= ?", until)
	}
	return db
}

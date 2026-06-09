// Package database 提供 PostgreSQL 数据库连接初始化功能。
package database

import (
	"fmt"

	"go_web/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres 根据配置连接 PostgreSQL，并设置连接池参数。
func NewPostgres(cfg config.PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
		cfg.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("open postgres handle: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}

// Package cache 提供 Redis 客户端初始化功能。
package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go_web/internal/config"
)

// NewRedis 根据配置创建 Redis 客户端，并验证连通性。
func NewRedis(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}

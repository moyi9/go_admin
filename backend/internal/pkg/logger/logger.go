// Package logger 封装 zap 日志库的初始化配置。
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go_web/internal/config"
)

// New 根据配置创建 zap.Logger，支持自定义日志级别和编码格式。
func New(cfg config.LogConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if cfg.Level != "" {
		if err := level.Set(cfg.Level); err != nil {
			return nil, fmt.Errorf("parse log level: %w", err)
		}
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.Encoding = cfg.Encoding
	if cfg.Encoding == "console" {
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	return zapCfg.Build()
}

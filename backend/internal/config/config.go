// Package config 管理应用配置，支持从 YAML 文件和环境变量加载。
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用顶层配置结构，包含所有子模块配置。
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Swagger  SwaggerConfig  `mapstructure:"swagger"`
	Seed     SeedConfig     `mapstructure:"seed"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

// AppConfig 应用基本信息。
type AppConfig struct {
	Name string `mapstructure:"name"` // 应用名称
	Env  string `mapstructure:"env"`  // 运行环境：development / production
}

// ServerConfig HTTP 服务器配置。
type ServerConfig struct {
	Host string `mapstructure:"host"` // 监听地址
	Port int    `mapstructure:"port"` // 监听端口
	Mode string `mapstructure:"mode"` // Gin 运行模式：debug / release / test
}

// PostgresConfig PostgreSQL 数据库连接配置。
type PostgresConfig struct {
	Host         string `mapstructure:"host"`          // 数据库主机地址
	Port         int    `mapstructure:"port"`          // 数据库端口
	User         string `mapstructure:"user"`          // 数据库用户
	Password     string `mapstructure:"password"`      // 数据库密码
	Database     string `mapstructure:"database"`      // 数据库名称
	SSLMode      string `mapstructure:"ssl_mode"`      // SSL 模式：disable / require
	Timezone     string `mapstructure:"timezone"`      // 时区，如 Asia/Shanghai
	MaxOpenConns int    `mapstructure:"max_open_conns"` // 最大打开连接数
	MaxIdleConns int    `mapstructure:"max_idle_conns"` // 最大空闲连接数
}

// RedisConfig Redis 缓存配置。
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`     // Redis 地址
	Password string `mapstructure:"password"` // Redis 密码
	DB       int    `mapstructure:"db"`       // Redis 数据库编号
}

// JWTConfig JWT 令牌配置。
type JWTConfig struct {
	Issuer         string        `mapstructure:"issuer"`           // 签发者
	Secret         string        `mapstructure:"secret"`           // 签名密钥
	AccessTokenTTL time.Duration `mapstructure:"access_token_ttl"` // Token 有效期
}

// LogConfig 日志配置。
type LogConfig struct {
	Level    string `mapstructure:"level"`    // 日志级别：debug / info / warn / error
	Encoding string `mapstructure:"encoding"` // 输出格式：console / json
}

// CORSConfig 跨域配置。
type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"` // 允许的域名
	AllowMethods []string `mapstructure:"allow_methods"` // 允许的 HTTP 方法
	AllowHeaders []string `mapstructure:"allow_headers"` // 允许的请求头
}

// SwaggerConfig Swagger 文档配置。
type SwaggerConfig struct {
	Enabled bool `mapstructure:"enabled"` // 是否启用 Swagger
}

// SeedConfig 种子数据配置，仅在开发环境启用。
type SeedConfig struct {
	Enabled       bool   `mapstructure:"enabled"`        // 是否启用自动初始化
	AdminUsername string `mapstructure:"admin_username"` // 默认管理员用户名
	AdminEmail    string `mapstructure:"admin_email"`    // 默认管理员邮箱
	AdminPassword string `mapstructure:"admin_password"` // 默认管理员密码
}

// UploadConfig 文件上传配置。
type UploadConfig struct {
	Dir         string   `mapstructure:"dir"`          // 文件存储根目录
	MaxSize     int64    `mapstructure:"max_size"`     // 单文件最大字节
	AllowedExts []string `mapstructure:"allowed_exts"` // 允许的扩展名
}

// Load 从 YAML 配置文件和环境变量加载完整配置，环境变量前缀 GO_WEB（如 GO_WEB_SERVER_PORT）。
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("GO_WEB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	setDefaults(v)

	if path != "" {
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置应用配置的默认值，避免因配置缺项导致启动失败。
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "go_web")
	v.SetDefault("app.env", "development")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("postgres.ssl_mode", "disable")
	v.SetDefault("postgres.timezone", "Asia/Shanghai")
	v.SetDefault("postgres.max_open_conns", 20)
	v.SetDefault("postgres.max_idle_conns", 10)
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("jwt.issuer", "go_web")
	v.SetDefault("jwt.access_token_ttl", "2h")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.encoding", "console")
	v.SetDefault("cors.allow_origins", []string{"*"})
	v.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allow_headers", []string{"Authorization", "Content-Type", "X-Request-ID"})
	v.SetDefault("swagger.enabled", true)
	v.SetDefault("seed.enabled", true)
	v.SetDefault("upload.dir", "./uploads")
	v.SetDefault("upload.max_size", 5<<20) // 5MB
	v.SetDefault("upload.allowed_exts", []string{".jpg", ".jpeg", ".png", ".gif", ".webp"})
}

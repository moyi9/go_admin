package bootstrap

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go_web/internal/config"
	"go_web/internal/modules/audit"
	"go_web/internal/modules/auth"
	"go_web/internal/modules/dashboard"
	"go_web/internal/modules/rbac"
	"go_web/internal/modules/notification"
	"go_web/internal/modules/upload"
	"go_web/internal/pkg/cache"
	"go_web/internal/pkg/database"
	"go_web/internal/pkg/logger"
	"go_web/internal/pkg/security"
	"go_web/internal/router"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB
	Redis  *redis.Client
	Router *gin.Engine
}

func New(ctx context.Context, configPath string, casbinModelPath string) (*App, error) {
	// 启动顺序：配置 -> 日志 -> 外部依赖 -> 数据表 -> 种子数据 -> 业务服务 -> 路由。
	// 这样可以保证服务真正监听端口前，数据库、Redis、权限策略都已经可用。
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	log, err := logger.New(cfg.Log)
	if err != nil {
		return nil, err
	}

	db, err := database.NewPostgres(cfg.Postgres)
	if err != nil {
		return nil, err
	}
	var redisClient *redis.Client
	redisClient, err = cache.NewRedis(ctx, cfg.Redis)
	if err != nil {
		log.Warn("redis unavailable, continuing without cache", zap.Error(err))
	} else {
		log.Info("redis connected", zap.String("addr", cfg.Redis.Addr))
	}
	if err := db.AutoMigrate(rbac.Models()...); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	if err := db.AutoMigrate(audit.Models()...); err != nil {
		return nil, fmt.Errorf("auto migrate audit: %w", err)
	}
	if err := db.AutoMigrate(notification.Models()...); err != nil {
		return nil, fmt.Errorf("auto migrate notification: %w", err)
	}
	if err := rbac.Seed(db, cfg.Seed); err != nil {
		return nil, fmt.Errorf("seed database: %w", err)
	}

	// Casbin 使用内存策略，启动时从 RBAC 数据表同步；业务授权变更后会重新同步。
	enforcer, err := casbin.NewEnforcer(casbinModelPath)
	if err != nil {
		return nil, fmt.Errorf("create casbin enforcer: %w", err)
	}

	jwtManager := security.NewJWTManager(cfg.JWT)
	rbacRepo := rbac.NewRepository(db)
	rbacService := rbac.NewService(rbacRepo, enforcer)
	if err := rbacService.SyncPolicies(); err != nil {
		return nil, err
	}

	auditRepo := audit.NewRepository(db)
	auditService := audit.NewService(auditRepo)

	dashboardRepo := dashboard.NewRepository(db)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	authService := auth.NewService(rbacRepo, jwtManager)
	notifRepo := notification.NewRepository(db)
	notifService := notification.NewService(notifRepo)
	notifHandler := notification.NewHandler(notifService, db)

	uploadHandler := upload.NewHandler(cfg.Upload, auditService)

	engine := router.New(router.Deps{
		Config:       cfg,
		Logger:       log,
		JWT:          jwtManager,
		Enforcer:     enforcer,
		AuthHandler:  auth.NewHandler(authService, auditService, notifService),
		RBACHandler:  rbac.NewHandler(rbacService, auditService),
		AuditHandler: audit.NewHandler(auditService),
		NotifHandler:   notifHandler,
		UploadHandler: uploadHandler,
		DashboardHandler: dashboardHandler,
	})

	return &App{
		Config: cfg,
		Logger: log,
		DB:     db,
		Redis:  redisClient,
		Router: engine,
	}, nil
}

// Close 释放 App 持有的外部资源（数据库连接、Redis 连接等）。
func (a *App) Close() {
	if a.Redis != nil {
		if err := a.Redis.Close(); err != nil {
			a.Logger.Error("close redis", zap.Error(err))
		}
	}
	if a.DB != nil {
		sqlDB, err := a.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				a.Logger.Error("close postgres", zap.Error(err))
			}
		}
	}
	a.Logger.Info("resources released")
}

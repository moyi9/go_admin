// Package router 注册所有 HTTP 路由和中间件，组织公开路由与需鉴权的业务路由。
package router

import (
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"go_web/internal/config"
	"go_web/internal/middleware"
	"go_web/internal/modules/audit"
	"go_web/internal/modules/auth"
	"go_web/internal/modules/dashboard"
	"go_web/internal/modules/rbac"
	"go_web/internal/modules/notification"
	"go_web/internal/modules/upload"
	"go_web/internal/pkg/response"
	"go_web/internal/pkg/security"
)

// Deps 路由依赖，用于解耦路由注册与具体模块初始化。
type Deps struct {
	Config        *config.Config
	Logger        *zap.Logger
	JWT           *security.JWTManager
	Enforcer      *casbin.Enforcer
	AuthHandler   *auth.Handler
	RBACHandler   *rbac.Handler
	AuditHandler  *audit.Handler
	UploadHandler *upload.Handler
	NotifHandler  *notification.Handler
	DashboardHandler *dashboard.Handler
}

// New 构建 Gin Router，注册全局中间件、健康检查、登录路由和受保护的业务路由。
func New(deps Deps) *gin.Engine {
	gin.SetMode(deps.Config.Server.Mode)
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(deps.Logger))
	r.Use(middleware.AccessLog(deps.Logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins: deps.Config.CORS.AllowOrigins,
		AllowMethods: deps.Config.CORS.AllowMethods,
		AllowHeaders: deps.Config.CORS.AllowHeaders,
		MaxAge:       12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})

	// 登录接口限流：每分钟每个 IP 最多 5 次请求，防止暴力破解。
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)

	// 公开接口 — 无需登录即可访问。
	api := r.Group("/api/v1")
	api.POST("/auth/login", loginLimiter.Middleware(), deps.AuthHandler.Login)

	// 受保护接口 — 需 JWT 鉴权 + Casbin 授权。
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(deps.JWT), middleware.Authorize(deps.Enforcer))
	protected.GET("/auth/me", deps.AuthHandler.Me)
	protected.PUT("/auth/password", deps.AuthHandler.UpdatePassword)
	protected.PUT("/auth/profile", deps.AuthHandler.UpdateProfile)

	protected.GET("/users", deps.RBACHandler.ListUsers)
	protected.POST("/users", deps.RBACHandler.CreateUser)
	protected.DELETE("/users/batch", deps.RBACHandler.BatchDeleteUsers)
	protected.GET("/users/:id", deps.RBACHandler.GetUser)
	protected.PUT("/users/:id", deps.RBACHandler.UpdateUser)
	protected.DELETE("/users/:id", deps.RBACHandler.DeleteUser)
	protected.POST("/users/:id/roles", deps.RBACHandler.AssignRoles)

	protected.GET("/roles", deps.RBACHandler.ListRoles)
	protected.POST("/roles", deps.RBACHandler.CreateRole)
	protected.DELETE("/roles/batch", deps.RBACHandler.BatchDeleteRoles)
	protected.PUT("/roles/:id", deps.RBACHandler.UpdateRole)
	protected.DELETE("/roles/:id", deps.RBACHandler.DeleteRole)
	protected.POST("/roles/:id/permissions", deps.RBACHandler.AssignPermissions)

	protected.GET("/permissions", deps.RBACHandler.ListPermissions)
	protected.POST("/permissions", deps.RBACHandler.CreatePermission)
	protected.DELETE("/permissions/batch", deps.RBACHandler.BatchDeletePermissions)
	protected.PUT("/permissions/:id", deps.RBACHandler.UpdatePermission)
	protected.DELETE("/permissions/:id", deps.RBACHandler.DeletePermission)

	protected.GET("/audit-logs", deps.AuditHandler.ListAuditLogs)

	protected.GET("/notifications", deps.NotifHandler.ListNotifications)
	protected.GET("/notifications/unread", deps.NotifHandler.CountUnread)
	protected.GET("/notifications/unread-list", deps.NotifHandler.ListUnreadNotifications)
	protected.POST("/notifications", deps.NotifHandler.SendNotification)
	protected.PUT("/notifications/:id/read", deps.NotifHandler.MarkAsRead)
	protected.PUT("/notifications/read-all", deps.NotifHandler.MarkAllAsRead)

	protected.GET("/dashboard/stats", deps.DashboardHandler.Stats)

	if deps.Config.Swagger.Enabled {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": 40400, "message": "route not found"})
	})

	return r
}

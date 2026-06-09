package rbac

import (
	"fmt"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// 每个测试独立的内存数据库，避免 cache=shared 导致的跨测试数据污染。
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(Models()...))
	return db
}

func seedTestUsers(t *testing.T, db *gorm.DB) {
	t.Helper()
	db.Create(&User{Username: "alice", Email: "alice@example.com", PasswordHash: "$2a$10$xxx", Status: UserStatusActive})
	db.Create(&User{Username: "bob", Email: "bob@example.com", PasswordHash: "$2a$10$xxx", Status: UserStatusActive})
	db.Create(&User{Username: "charlie", Email: "charlie@example.com", PasswordHash: "$2a$10$xxx", Status: UserStatusDisabled})
}

func seedTestRolesAndPerms(t *testing.T, db *gorm.DB) {
	t.Helper()
	db.Create(&Role{Code: "test_admin", Name: "Admin"})
	db.Create(&Role{Code: "test_editor", Name: "Editor"})
	db.Create(&Role{Code: "test_viewer", Name: "Viewer"})
	db.Create(&Permission{Code: "test.get.users", Name: "List users", Method: "GET", Path: "/api/v1/users"})
	db.Create(&Permission{Code: "test.post.users", Name: "Create user", Method: "POST", Path: "/api/v1/users"})
	db.Create(&Permission{Code: "test.get.roles", Name: "List roles", Method: "GET", Path: "/api/v1/roles"})
}

func TestListUsersPaged(t *testing.T) {
	db := newTestDB(t)
	seedTestUsers(t, db)

	enforcer, _ := casbin.NewEnforcer("../../../configs/casbin_model.conf")
	repo := NewRepository(db)
	svc := NewService(repo, enforcer)

	users, total, err := svc.ListUsersPaged(0, 2, ListUsersQuery{})
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, users, 2)

	users, total, err = svc.ListUsersPaged(2, 2, ListUsersQuery{})
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, users, 1)
}

func TestListRolesPaged(t *testing.T) {
	db := newTestDB(t)
	seedTestRolesAndPerms(t, db)

	enforcer, _ := casbin.NewEnforcer("../../../configs/casbin_model.conf")
	repo := NewRepository(db)
	svc := NewService(repo, enforcer)

	roles, total, err := svc.ListRolesPaged(0, 2, ListRolesQuery{})
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, roles, 2)
}

func TestListPermissionsPaged(t *testing.T) {
	db := newTestDB(t)
	seedTestRolesAndPerms(t, db)

	enforcer, _ := casbin.NewEnforcer("../../../configs/casbin_model.conf")
	repo := NewRepository(db)
	svc := NewService(repo, enforcer)

	perms, total, err := svc.ListPermissionsPaged(0, 2, ListPermissionsQuery{})
	require.NoError(t, err)
	require.Equal(t, int64(3), total)
	require.Len(t, perms, 2)
}

func TestCreateUserDefaultStatus(t *testing.T) {
	db := newTestDB(t)
	enforcer, _ := casbin.NewEnforcer("../../../configs/casbin_model.conf")
	svc := NewService(NewRepository(db), enforcer)

	user, err := svc.CreateUser(CreateUserRequest{
		Username: "dave",
		Email:    "dave@example.com",
		Password: "password123456",
	})
	require.NoError(t, err)
	require.Equal(t, UserStatusActive, user.Status)
	require.NotEmpty(t, user.PasswordHash)
	require.NotEqual(t, "password123456", user.PasswordHash)
}

func TestCreateUserDuplicate(t *testing.T) {
	db := newTestDB(t)
	enforcer, _ := casbin.NewEnforcer("../../../configs/casbin_model.conf")
	svc := NewService(NewRepository(db), enforcer)

	_, err := svc.CreateUser(CreateUserRequest{
		Username: "eve_unique",
		Email:    "eve@example.com",
		Password: "password123456",
	})
	require.NoError(t, err)

	_, err = svc.CreateUser(CreateUserRequest{
		Username: "eve_unique",
		Email:    "eve2@example.com",
		Password: "password123456",
	})
	require.Error(t, err)
}

func TestAssignRolesValidatesAllExist(t *testing.T) {
	db := newTestDB(t)
	enforcer, _ := casbin.NewEnforcer("../../../configs/casbin_model.conf")
	svc := NewService(NewRepository(db), enforcer)

	db.Create(&User{Username: "frank_u", Email: "frank_u@example.com", PasswordHash: "x", Status: UserStatusActive})
	db.Create(&Role{Code: "tester_u", Name: "Tester"})

	var user User
	db.First(&user, "username = ?", "frank_u")
	err := svc.AssignRoles(user.ID, []uint{999}) // 999 不存在
	require.Error(t, err)
	require.Contains(t, err.Error(), "some roles do not exist")
}

func TestSyncPolicies(t *testing.T) {
	db := newTestDB(t)
	enforcer, _ := casbin.NewEnforcer("../../../configs/casbin_model.conf")
	svc := NewService(NewRepository(db), enforcer)

	// 创建角色 + 权限 + 建立关联
	role := Role{Code: "qa_sync", Name: "QA Sync"}
	require.NoError(t, db.Create(&role).Error)
	perm := Permission{Code: "get.healthz.sync", Name: "Health Sync", Method: "GET", Path: "/healthz"}
	require.NoError(t, db.Create(&perm).Error)
	require.NoError(t, db.Model(&role).Association("Permissions").Append(&perm))

	// 创建用户 + 分配角色
	user := User{Username: "grace_sync", Email: "grace_s@example.com", PasswordHash: "x", Status: UserStatusActive}
	require.NoError(t, db.Create(&user).Error)
	require.NoError(t, db.Model(&user).Association("Roles").Append(&role))

	// 同步策略
	require.NoError(t, svc.SyncPolicies())

	// 验证 Casbin 策略：用户 ID → 角色 code → 权限 path + method
	allowed, err := enforcer.Enforce(fmt.Sprintf("%d", user.ID), "/healthz", "GET")
	require.NoError(t, err)
	require.True(t, allowed, "user %d should have access to GET /healthz", user.ID)
}

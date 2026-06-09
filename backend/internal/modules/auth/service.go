package auth

import (
	"go_web/internal/modules/rbac"
	"go_web/internal/pkg/apperror"
	"go_web/internal/pkg/security"
)

// Service 处理用户认证业务逻辑。
type Service struct {
	repo *rbac.Repository
	jwt  *security.JWTManager
}

// NewService 创建认证 Service。
func NewService(repo *rbac.Repository, jwt *security.JWTManager) *Service {
	return &Service{repo: repo, jwt: jwt}
}

// Login 验证用户名密码，校验用户状态，生成 JWT Token 并返回。
func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Status != rbac.UserStatusActive || !security.CheckPassword(user.PasswordHash, req.Password) {
		return nil, apperror.New(apperror.CodeUnauthorized, "invalid username or password")
	}

	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, role.Code)
	}
	token, err := s.jwt.Generate(user.ID, user.Username, roles)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		User:        *user,
	}, nil
}

// UpdatePassword 修改当前用户密码。验证旧密码正确后，加密新密码并更新数据库。
func (s *Service) UpdatePassword(userID uint, req UpdatePasswordRequest) error {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return apperror.New(apperror.CodeNotFound, "user not found")
	}
	if !security.CheckPassword(user.PasswordHash, req.CurrentPassword) {
		return apperror.New(apperror.CodeInvalidArgument, "current password is incorrect")
	}
	hashed, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserPassword(userID, hashed)
}

// UpdateProfile 更新当前用户个人信息（昵称、邮箱、手机号、性别）。
func (s *Service) UpdateProfile(userID uint, req UpdateProfileRequest) (*rbac.User, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.New(apperror.CodeNotFound, "user not found")
	}
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Gender = req.Gender
	user.AvatarURL = req.AvatarURL
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Me 根据用户 ID 查询用户详情。
func (s *Service) Me(userID uint) (*rbac.User, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.New(apperror.CodeNotFound, "user not found")
	}
	return user, nil
}

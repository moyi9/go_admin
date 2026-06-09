package notification

import "go_web/internal/pkg/apperror"

// Service 处理通知业务逻辑。
type Service struct {
	repo *Repository
}

// NewService 创建通知 Service。
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Send 给指定用户发送一条通知。
func (s *Service) Send(userID uint, notifType, title, content, link string) error {
	return s.repo.Create(&Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Content: content,
		Link:    link,
	})
}

// ListResponse 通知分页查询响应。
type ListResponse struct {
	List  []Notification `json:"list"`
	Total int64          `json:"total"`
}

// List 分页查询当前用户的通知。
func (s *Service) List(userID uint, offset, limit int) (*ListResponse, error) {
	list, total, err := s.repo.ListByUser(userID, offset, limit)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "list notifications failed", err)
	}
	return &ListResponse{List: list, Total: total}, nil
}

// ListUnread 查询最新 N 条未读通知（供铃铛下拉预览用）。
func (s *Service) ListUnread(userID uint, limit int) ([]Notification, error) {
	return s.repo.ListUnread(userID, limit)
}

// CountUnread 返回当前用户的未读通知数。
func (s *Service) CountUnread(userID uint) (int64, error) {
	return s.repo.CountUnread(userID)
}

// MarkAsRead 将指定用户的一条通知标记为已读。
func (s *Service) MarkAsRead(userID, id uint) error {
	return s.repo.MarkAsRead(userID, id)
}

// MarkAllAsRead 将指定用户的全部通知标记为已读。
func (s *Service) MarkAllAsRead(userID uint) error {
	return s.repo.MarkAllAsRead(userID)
}

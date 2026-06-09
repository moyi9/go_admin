package audit

import "go_web/internal/pkg/apperror"

// Service 处理操作日志业务逻辑。
type Service struct {
	repo *Repository
}

// NewService 创建操作日志 Service。
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Log 写入一条操作日志。
func (s *Service) Log(log *AuditLog) error {
	return s.repo.Insert(log)
}

// ListAuditLogsQuery 日志列表查询参数。
type ListQuery struct {
	Action   string
	Resource string
	Keyword  string
	Since    string
	Until    string
}

// ListResponse 日志分页查询响应。
type ListResponse struct {
	List  []AuditLog `json:"list"`
	Total int64      `json:"total"`
}

// List 分页查询操作日志。
func (s *Service) List(offset, limit int, q ListQuery) (*ListResponse, error) {
	repoQ := ListAuditLogsQuery{
		Action:   q.Action,
		Resource: q.Resource,
		Keyword:  q.Keyword,
	}
	total, err := s.repo.Count(repoQ, q.Since, q.Until)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "count audit logs failed", err)
	}
	list, err := s.repo.ListPaged(offset, limit, repoQ, q.Since, q.Until)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "list audit logs failed", err)
	}
	if list == nil {
		list = []AuditLog{}
	}
	return &ListResponse{List: list, Total: total}, nil
}

package dashboard

// Service 封装仪表盘统计业务逻辑。
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Stats 获取全部仪表盘聚合统计。
func (s *Service) Stats() (*Stats, error) {
	gender, err := s.repo.GenderDistribution()
	if err != nil {
		return nil, err
	}
	status, err := s.repo.StatusDistribution()
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.RoleDistribution()
	if err != nil {
		return nil, err
	}
	return &Stats{
		GenderDistribution: gender,
		StatusDistribution: status,
		RoleDistribution:   roles,
	}, nil
}

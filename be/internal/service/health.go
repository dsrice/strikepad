package service

type HealthService interface {
	Check() map[string]string
}

type healthService struct{}

func NewHealthService() HealthService {
	return &healthService{}
}

func (s *healthService) Check() map[string]string {
	return map[string]string{
		"status": "ok",
	}
}

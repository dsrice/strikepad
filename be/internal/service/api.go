package service

type APIService interface {
	GetTestMessage() map[string]string
}

type apiService struct{}

func NewAPIService() APIService {
	return &apiService{}
}

func (s *apiService) GetTestMessage() map[string]string {
	return map[string]string{
		"message": "API endpoint working",
	}
}

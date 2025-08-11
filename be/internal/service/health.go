package service

import "strikepad-backend/internal/dto"

type healthService struct{}

func NewHealthService() HealthServiceInterface {
	return &healthService{}
}

func (s *healthService) GetHealth() *dto.HealthResponse {
	return &dto.HealthResponse{
		Status:  "ok",
		Message: "Server is healthy",
	}
}

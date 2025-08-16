package handler

import (
	"net/http"

	"strikepad-backend/internal/service"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	healthService service.HealthServiceInterface
}

func NewHealthHandler(healthService service.HealthServiceInterface) HealthHandlerInterface {
	return &HealthHandler{
		healthService: healthService,
	}
}

// Check handles health check endpoint
// @Summary Health check
// @Description Check if the service is healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthResponse "Service is healthy"
// @Router /health [get]
func (h *HealthHandler) Check(c echo.Context) error {
	result := h.healthService.GetHealth()
	return c.JSON(http.StatusOK, result)
}

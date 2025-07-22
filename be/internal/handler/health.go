package handler

import (
	"net/http"

	"strikepad-backend/internal/service"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	healthService service.HealthService
}

func NewHealthHandler(healthService service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

func (h *HealthHandler) Health(c echo.Context) error {
	result := h.healthService.Check()
	return c.JSON(http.StatusOK, result)
}

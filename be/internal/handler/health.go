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

func (h *HealthHandler) Check(c echo.Context) error {
	result := h.healthService.GetHealth()
	return c.JSON(http.StatusOK, result)
}

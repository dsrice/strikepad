package handler

import (
	"net/http"

	"strikepad-backend/internal/service"

	"github.com/labstack/echo/v4"
)

type APIHandler struct {
	apiService service.APIService
}

func NewAPIHandler(apiService service.APIService) *APIHandler {
	return &APIHandler{
		apiService: apiService,
	}
}

func (h *APIHandler) Test(c echo.Context) error {
	result := h.apiService.GetTestMessage()
	return c.JSON(http.StatusOK, result)
}

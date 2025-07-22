package main

import (
	"log"
	"net/http"

	"strikepad-backend/internal/container"
	"strikepad-backend/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	c := container.BuildContainer()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from StrikePad Backend!")
	})

	err := c.Invoke(func(healthHandler *handler.HealthHandler, apiHandler *handler.APIHandler) {
		e.GET("/health", healthHandler.Health)
		e.GET("/api/test", apiHandler.Test)
	})

	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}

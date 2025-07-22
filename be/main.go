package main

import (
	"net/http"

	"github.com/la
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from StrikePad Backend!")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	e.GET("/api/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "API endpoint working",
		})

	fmt.Println("Server starting on port 8080...")
	e.Logger.Fatal(e.Start(":8080"))
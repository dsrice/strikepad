package main

import (
	"net/http"
	"strikepad-manage-tool/handlers"
	"strikepad-manage-tool/templates"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echoインスタンス作成
	e := echo.New()

	// テンプレートレンダラー設定
	e.Renderer = templates.NewTemplateRenderer()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 認証ハンドラー初期化
	authHandler := handlers.NewAuthHandler()

	// 認証関連ルート
	e.GET("/login", authHandler.ShowLogin)
	e.POST("/login", authHandler.Login)
	e.GET("/logout", authHandler.Logout)
	e.GET("/dashboard", authHandler.ShowDashboard)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/login")
	})

	// 管理用API群
	api := e.Group("/api")
	api.GET("/health", healthCheck)

	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}

// ヘルスチェックエンドポイント
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "strikepad-manage-tool",
	})
}
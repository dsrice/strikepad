package main

import (
	"log"
	"net/http"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/handlers"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/templates"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// データベース接続
	dbConfig := config.NewDatabaseConfig()
	db, err := config.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(db)

	// Echoインスタンス作成
	e := echo.New()

	// テンプレートレンダラー設定
	e.Renderer = templates.NewTemplateRenderer()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ハンドラー初期化
	authHandler := handlers.NewAuthHandler()
	adminHandler := handlers.NewAdminHandler(userRepo)

	// 認証関連ルート
	e.GET("/login", authHandler.ShowLogin)
	e.POST("/login", authHandler.Login)
	e.GET("/logout", authHandler.Logout)
	e.GET("/dashboard", authHandler.ShowDashboard)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/login")
	})

	// 管理機能ルート
	admin := e.Group("/admin")
	admin.GET("/users", adminHandler.ShowUsers)
	admin.GET("/users/:id", adminHandler.ShowUserDetail)
	admin.POST("/users/:id/status", adminHandler.UpdateUserStatus)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)

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

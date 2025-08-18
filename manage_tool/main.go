package main

import (
	"log"
	"net/http"
	"os"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/handlers"
	"strikepad-manage-tool/middleware"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/templates"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
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
	adminRepo := repository.NewAdminRepository(db)
	makerRepo := repository.NewMakerRepository(db)

	// Echoインスタンス作成
	e := echo.New()

	// テンプレートレンダラー設定
	e.Renderer = templates.NewTemplateRenderer()

	// セッション設定
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "your-secret-key-change-in-production" // 本番環境では必ず変更
		log.Println("Warning: Using default session secret. Set SESSION_SECRET environment variable in production.")
	}
	sessionStore := middleware.CreateSessionStore(sessionSecret)
	e.Use(session.Middleware(sessionStore))

	// ミドルウェア設定
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// ハンドラー初期化
	authHandler := handlers.NewAuthHandler(adminRepo)
	adminHandler := handlers.NewAdminHandler(userRepo)
	makerHandler := handlers.NewMakerHandler(makerRepo)

	// 認証が不要なルート
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/login")
	})
	e.GET("/login", authHandler.ShowLogin)
	e.POST("/auth/login", authHandler.Login)

	// 管理用API（認証不要）
	api := e.Group("/api")
	api.GET("/health", healthCheck)

	// セッション認証が必要なルート
	protected := e.Group("")
	protected.Use(middleware.SessionAuth())

	// 認証が必要なルート
	protected.GET("/dashboard", adminHandler.ShowDashboard)
	protected.GET("/logout", authHandler.Logout)

	// 管理機能ルート（認証必要）
	admin := protected.Group("/admin")

	// ユーザー管理
	admin.GET("/users", adminHandler.ShowUsers)
	admin.GET("/users/:id", adminHandler.ShowUserDetail)
	admin.POST("/users/:id/status", adminHandler.UpdateUserStatus)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)

	// メーカー管理
	admin.GET("/makers", makerHandler.ShowMakers)
	admin.GET("/makers/create", makerHandler.ShowCreateMaker)
	admin.POST("/makers/create", makerHandler.CreateMaker)
	admin.GET("/makers/:id", makerHandler.ShowMakerDetail)
	admin.DELETE("/makers/:id", makerHandler.DeleteMaker)
	admin.GET("/makers/stats", makerHandler.GetMakerStats)

	log.Println("Starting Strikepad Management Tool on :8080")
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
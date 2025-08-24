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

	// S3クライアント接続
	s3Client, err := config.NewS3Client()
	if err != nil {
		log.Fatal("Failed to connect to S3:", err)
	}

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	makerRepo := repository.NewMakerRepository(db)
	coreRepo := repository.NewCoreRepository(db)
	coverRepo := repository.NewCoverRepository(db)

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
	makerHandler := handlers.NewMakerHandler(makerRepo, s3Client)
	coreHandler := handlers.NewCoreHandler(coreRepo, makerRepo)
	coverHandler := handlers.NewCoverHandler(coverRepo, makerRepo)

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
	admin.GET("/makers/:id/edit", makerHandler.ShowEditMaker)
	admin.POST("/makers/:id/edit", makerHandler.UpdateMaker)
	admin.DELETE("/makers/:id", makerHandler.DeleteMaker)
	admin.GET("/makers/stats", makerHandler.GetMakerStats)
	admin.GET("/makers/:id/logo-url", makerHandler.GetLogoPresignedURL)
	admin.POST("/makers/:id/upload-url", makerHandler.GetUploadPresignedURL)
	admin.POST("/makers/:id/confirm-upload", makerHandler.ConfirmLogoUpload)

	// コア管理
	admin.GET("/cores", coreHandler.ShowCores)
	admin.GET("/cores/create", coreHandler.ShowCreateCore)
	admin.POST("/cores/create", coreHandler.CreateCore)
	admin.GET("/cores/:id", coreHandler.ShowCoreDetail)
	admin.DELETE("/cores/:id", coreHandler.DeleteCore)
	admin.GET("/cores/stats", coreHandler.GetCoreStats)

	// カバー管理
	admin.GET("/covers", coverHandler.ShowCovers)
	admin.GET("/covers/create", coverHandler.ShowCreateCover)
	admin.POST("/covers", coverHandler.CreateCover)
	admin.GET("/covers/stats", coverHandler.GetCoverStats)
	admin.GET("/covers/:id", coverHandler.ShowCoverDetail)
	admin.GET("/covers/:id/edit", coverHandler.ShowEditCover)
	admin.POST("/covers/:id/update", coverHandler.UpdateCover)
	admin.DELETE("/covers/:id", coverHandler.DeleteCover)

	log.Println("Starting Strikepad Management Tool on :8081")
	// サーバー起動
	e.Logger.Fatal(e.Start(":8081"))
}

// ヘルスチェックエンドポイント
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "strikepad-manage-tool",
	})
}

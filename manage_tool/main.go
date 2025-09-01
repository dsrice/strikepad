package main

import (
	"log"
	"net/http"
	"os"

	"strikepad-manage-tool/container"
	"strikepad-manage-tool/middleware"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// DIコンテナー初期化
	c, err := container.BuildContainer()
	if err != nil {
		log.Fatal("Failed to build DI container:", err)
	}

	// Echoインスタンス作成
	e := echo.New()

	// テンプレートレンダラー設定
	err = c.Invoke(func(renderer echo.Renderer) {
		e.Renderer = renderer
	})
	if err != nil {
		log.Fatal("Failed to set template renderer:", err)
	}

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

	// ハンドラー初期化とルート設定
	err = c.Invoke(setupRoutes(e))
	if err != nil {
		log.Fatal("Failed to setup routes:", err)
	}

	log.Println("Starting Strikepad Management Tool on :8081")
	// サーバー起動
	e.Logger.Fatal(e.Start(":8081"))
}

// setupRoutes は依存関係注入されたハンドラーでルートを設定する関数を返す
func setupRoutes(e *echo.Echo) interface{} {
	return func(handlers *container.HandlersContainer) error {
		// 認証が不要なルート
		e.GET("/", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, "/login")
		})
		e.GET("/login", handlers.AuthHandler.ShowLogin)
		e.POST("/auth/login", handlers.AuthHandler.Login)

		// 管理用API（認証不要）
		api := e.Group("/api")
		api.GET("/health", healthCheck)

		// セッション認証が必要なルート
		protected := e.Group("")
		protected.Use(middleware.SessionAuth())

		// 認証が必要なルート
		protected.GET("/dashboard", handlers.AdminHandler.ShowDashboard)
		protected.GET("/logout", handlers.AuthHandler.Logout)

		// 管理機能ルート（認証必要）
		admin := protected.Group("/admin")

		// ユーザー管理
		admin.GET("/users", handlers.AdminHandler.ShowUsers)
		admin.GET("/users/:id", handlers.AdminHandler.ShowUserDetail)
		admin.POST("/users/:id/status", handlers.AdminHandler.UpdateUserStatus)
		admin.DELETE("/users/:id", handlers.AdminHandler.DeleteUser)

		// メーカー管理
		admin.GET("/makers", handlers.MakerHandler.ShowMakers)
		admin.GET("/makers/create", handlers.MakerHandler.ShowCreateMaker)
		admin.POST("/makers/create", handlers.MakerHandler.CreateMaker)
		admin.GET("/makers/:id", handlers.MakerHandler.ShowMakerDetail)
		admin.GET("/makers/:id/edit", handlers.MakerHandler.ShowEditMaker)
		admin.POST("/makers/:id/edit", handlers.MakerHandler.UpdateMaker)
		admin.DELETE("/makers/:id", handlers.MakerHandler.DeleteMaker)
		admin.GET("/makers/stats", handlers.MakerHandler.GetMakerStats)
		admin.GET("/makers/:id/logo-url", handlers.MakerHandler.GetLogoPresignedURL)
		admin.POST("/makers/:id/upload-url", handlers.MakerHandler.GetUploadPresignedURL)
		admin.POST("/makers/:id/confirm-upload", handlers.MakerHandler.ConfirmLogoUpload)

		// コア管理
		admin.GET("/cores", handlers.CoreHandler.ShowCores)
		admin.GET("/cores/create", handlers.CoreHandler.ShowCreateCore)
		admin.POST("/cores/create", handlers.CoreHandler.CreateCore)
		admin.GET("/cores/:id", handlers.CoreHandler.ShowCoreDetail)
		admin.GET("/cores/:id/edit", handlers.CoreHandler.ShowEditCore)
		admin.POST("/cores/:id/edit", handlers.CoreHandler.UpdateCore)
		admin.DELETE("/cores/:id", handlers.CoreHandler.DeleteCore)
		admin.GET("/cores/stats", handlers.CoreHandler.GetCoreStats)

		// カバー管理
		admin.GET("/covers", handlers.CoverHandler.ShowCovers)
		admin.GET("/covers/create", handlers.CoverHandler.ShowCreateCover)
		admin.POST("/covers", handlers.CoverHandler.CreateCover)
		admin.GET("/covers/stats", handlers.CoverHandler.GetCoverStats)
		admin.GET("/covers/:id", handlers.CoverHandler.ShowCoverDetail)
		admin.GET("/covers/:id/edit", handlers.CoverHandler.ShowEditCover)
		admin.POST("/covers/:id/update", handlers.CoverHandler.UpdateCover)
		admin.DELETE("/covers/:id", handlers.CoverHandler.DeleteCover)

		// ボール管理
		admin.GET("/balls", handlers.BallHandler.ShowBalls)
		admin.GET("/balls/create", handlers.BallHandler.ShowCreateBall)
		admin.POST("/balls/create", handlers.BallHandler.CreateBall)
		admin.GET("/balls/:id", handlers.BallHandler.ShowBallDetail)
		admin.GET("/balls/:id/edit", handlers.BallHandler.ShowEditBall)
		admin.POST("/balls/:id/edit", handlers.BallHandler.UpdateBall)
		admin.DELETE("/balls/:id", handlers.BallHandler.DeleteBall)
		admin.GET("/balls/stats", handlers.BallHandler.GetBallStats)

		// API エンドポイント
		e.GET("/admin/api/makers/:id/cores-and-covers", handlers.BallHandler.GetCoresAndCoversByMaker)

		return nil
	}
}

// ヘルスチェックエンドポイント
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "strikepad-manage-tool",
	})
}
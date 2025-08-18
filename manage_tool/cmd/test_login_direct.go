package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

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
	fmt.Println("=== 直接ログインテスト ===")

	// データベース接続
	dbConfig := config.NewDatabaseConfig()
	db, err := config.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	// リポジトリ初期化
	adminRepo := repository.NewAdminRepository(db)

	// Echoインスタンス作成
	e := echo.New()

	// テンプレートレンダラー設定
	e.Renderer = templates.NewTemplateRenderer()

	// セッション設定
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "test-secret-key-for-debugging"
	}
	sessionStore := middleware.CreateSessionStore(sessionSecret)
	e.Use(session.Middleware(sessionStore))

	// ミドルウェア設定
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	// ハンドラー初期化
	authHandler := handlers.NewAuthHandler(adminRepo)

	// ルート設定
	e.GET("/login", authHandler.ShowLogin)
	e.POST("/auth/login", authHandler.Login)

	fmt.Println("1. ログインページにアクセステスト")

	// ログインページのテスト
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := authHandler.ShowLogin(c); err != nil {
		log.Printf("ログインページ表示エラー: %v", err)
	} else {
		fmt.Printf("✅ ログインページ表示成功 (ステータス: %d)\n", rec.Code)
	}

	fmt.Println("\n2. ログイン処理テスト")

	// ログイン処理のテスト
	form := url.Values{}
	form.Add("login_id", "admin")
	form.Add("password", "admin123")

	req = httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()

	fmt.Printf("テストログイン: login_id=admin, password=admin123\n")

	// Echo のルーター経由でリクエストを実行（ミドルウェアが実行される）
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusFound {
		fmt.Printf("❌ ログイン処理に失敗しました (ステータス: %d)\n", rec.Code)
	} else {
		fmt.Printf("✅ ログイン処理完了 (ステータス: %d)\n", rec.Code)
		fmt.Printf("レスポンスヘッダー:\n")
		for name, values := range rec.Header() {
			for _, value := range values {
				fmt.Printf("  %s: %s\n", name, value)
			}
		}

		// リダイレクト確認
		if rec.Code == http.StatusFound {
			location := rec.Header().Get("Location")
			fmt.Printf("✅ リダイレクト先: %s\n", location)
		}

		// セッション情報の確認
		fmt.Printf("レスポンス本文の最初の100文字:\n%s...\n", string(rec.Body.Bytes()[:min(100, len(rec.Body.Bytes()))]))
	}

	fmt.Println("\n✅ 直接ログインテスト完了")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
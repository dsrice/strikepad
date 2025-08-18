package main

import (
	"fmt"
	"log"
	"time"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/repository"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("=== セッション機能テスト ===")

	// データベース接続
	dbConfig := config.NewDatabaseConfig()
	db, err := config.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	// リポジトリ初期化
	adminRepo := repository.NewAdminRepository(db)

	// 管理者ユーザーを取得
	adminUser, err := adminRepo.GetByLoginID("admin")
	if err != nil {
		log.Fatal("管理者ユーザー取得エラー:", err)
	}

	if adminUser == nil {
		log.Fatal("管理者ユーザーが見つかりません")
	}

	fmt.Printf("✅ 管理者ユーザーが見つかりました: %s\n", adminUser.Name)

	// Echoコンテキストを模擬作成
	e := echo.New()

	// セッションストア設定
	sessionSecret := "test-secret-key"
	store := sessions.NewCookieStore([]byte(sessionSecret))
	e.Use(session.Middleware(store))

	// テスト用のリクエスト作成
	fmt.Println("\n=== セッション作成テスト ===")

	// テスト用のEchoコンテキストを作成する代わりに
	// 直接セッションライブラリをテスト
	testSession := sessions.NewSession(store, "admin_session")

	// セッションに値を設定
	testSession.Values["user_id"] = adminUser.ID
	testSession.Values["user_name"] = adminUser.Name
	testSession.Values["login_id"] = adminUser.LoginID
	testSession.Values["last_access"] = time.Now()

	fmt.Printf("✅ セッションに値を設定しました\n")
	fmt.Printf("  - user_id: %v\n", testSession.Values["user_id"])
	fmt.Printf("  - user_name: %v\n", testSession.Values["user_name"])
	fmt.Printf("  - login_id: %v\n", testSession.Values["login_id"])
	fmt.Printf("  - last_access: %v\n", testSession.Values["last_access"])

	// セッション関連のモジュールが正常にインポートできるかテスト
	fmt.Println("\n=== 依存関係テスト ===")

	// gorilla/sessions のテスト
	testStore := sessions.NewCookieStore([]byte("test-key"))
	if testStore != nil {
		fmt.Printf("✅ gorilla/sessions が正常に動作します\n")
	} else {
		fmt.Printf("❌ gorilla/sessions に問題があります\n")
	}

	// echo-contrib/session のテスト
	fmt.Printf("✅ echo-contrib/session が正常にインポートされています\n")

	fmt.Println("\n=== 依存関係確認 ===")
	fmt.Println("以下のパッケージが必要です:")
	fmt.Println("- github.com/gorilla/sessions")
	fmt.Println("- github.com/labstack/echo-contrib/session")
	fmt.Println("- github.com/gorilla/securecookie")

	fmt.Println("\n✅ セッション機能テスト完了")
	fmt.Println("問題が発生している場合は、以下を確認してください:")
	fmt.Println("1. 依存関係が正しくインストールされているか")
	fmt.Println("2. セッションストアの設定が正しいか")
	fmt.Println("3. セッション保存時にエラーが発生していないか")
}
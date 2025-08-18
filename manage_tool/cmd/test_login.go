package main

import (
	"fmt"
	"log"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/utils"
)

func main() {
	fmt.Println("=== 初期管理ユーザーのログインテスト ===")

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

	fmt.Printf("✅ 管理者ユーザーが見つかりました\n")
	fmt.Printf("ID: %d\n", adminUser.ID)
	fmt.Printf("ログインID: %s\n", adminUser.LoginID)
	fmt.Printf("表示名: %s\n", adminUser.Name)

	// パスワードテスト
	fmt.Println("\n=== パスワードテスト ===")
	testPassword := "admin123"

	if utils.CheckPassword(testPassword, adminUser.Password) {
		fmt.Printf("✅ パスワード 'admin123' が正しく検証されました\n")
	} else {
		fmt.Printf("❌ パスワード検証に失敗しました\n")
	}

	// 間違ったパスワードのテスト
	wrongPassword := "wrongpassword"
	if !utils.CheckPassword(wrongPassword, adminUser.Password) {
		fmt.Printf("✅ 間違ったパスワード '%s' が正しく拒否されました\n", wrongPassword)
	} else {
		fmt.Printf("❌ 間違ったパスワードが受け入れられました（問題あり）\n")
	}

	fmt.Println("\n✅ 初期管理ユーザーのテスト完了")
	fmt.Println("管理ツールにログインできます:")
	fmt.Println("- URL: http://localhost:8080/login")
	fmt.Println("- ログインID: admin")
	fmt.Println("- パスワード: admin123")
}
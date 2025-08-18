//go:build tools
// +build tools

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/utils"

	"golang.org/x/term"
)

func main() {
	fmt.Println("=== Strikepad 管理ツール セットアップ ===")
	fmt.Println("初期管理者ユーザーを作成します\n")

	// データベース接続
	dbConfig := config.NewDatabaseConfig()
	db, err := config.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	// リポジトリ初期化
	adminRepo := repository.NewAdminRepository(db)

	// ユーザー入力
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("ログインID: ")
	loginID, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("入力エラー:", err)
	}
	loginID = strings.TrimSpace(loginID)

	if len(loginID) == 0 {
		log.Fatal("ログインIDは必須です")
	}

	// 既存チェック
	existingUser, err := adminRepo.GetByLoginID(loginID)
	if err != nil {
		log.Fatal("データベースエラー:", err)
	}
	if existingUser != nil {
		log.Fatal("このログインIDは既に使用されています")
	}

	fmt.Print("表示名: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("入力エラー:", err)
	}
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		log.Fatal("表示名は必須です")
	}

	// パスワード入力（隠す）
	fmt.Print("パスワード: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("パスワード入力エラー:", err)
	}
	password := string(passwordBytes)
	fmt.Println() // 改行

	if len(password) < 6 {
		log.Fatal("パスワードは6文字以上で設定してください")
	}

	fmt.Print("パスワード（確認）: ")
	confirmPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("パスワード確認入力エラー:", err)
	}
	confirmPassword := string(confirmPasswordBytes)
	fmt.Println() // 改行

	if password != confirmPassword {
		log.Fatal("パスワードが一致しません")
	}

	// パスワードハッシュ化
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal("パスワードハッシュ化エラー:", err)
	}

	// 管理者ユーザー作成
	adminUser := &models.AdminUser{
		LoginID:   loginID,
		Password:  hashedPassword,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	err = adminRepo.Create(adminUser)
	if err != nil {
		log.Fatal("管理者ユーザー作成エラー:", err)
	}

	fmt.Printf("\n✅ 管理者ユーザーを作成しました\n")
	fmt.Printf("ログインID: %s\n", loginID)
	fmt.Printf("表示名: %s\n", name)
	fmt.Printf("ID: %d\n", adminUser.ID)
	fmt.Println("\n管理ツールにログインできます: http://localhost:8080/login")
}

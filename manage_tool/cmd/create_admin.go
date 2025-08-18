package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"strikepad-manage-tool/config"
	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository"
	"strikepad-manage-tool/utils"
)

func main() {
	loginID := flag.String("login", "", "ログインID (必須)")
	name := flag.String("name", "", "表示名 (必須)")
	password := flag.String("password", "", "パスワード (必須)")
	flag.Parse()

	if *loginID == "" || *name == "" || *password == "" {
		fmt.Println("使用方法:")
		fmt.Println("go run cmd/create_admin.go -login=admin -name=Administrator -password=yourpassword")
		fmt.Println("")
		flag.Usage()
		return
	}

	if len(*password) < 6 {
		log.Fatal("パスワードは6文字以上で設定してください")
	}

	fmt.Println("=== 管理者ユーザー作成 ===")

	// データベース接続
	dbConfig := config.NewDatabaseConfig()
	db, err := config.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	// リポジトリ初期化
	adminRepo := repository.NewAdminRepository(db)

	// 既存チェック
	existingUser, err := adminRepo.GetByLoginID(*loginID)
	if err != nil {
		log.Fatal("データベースエラー:", err)
	}
	if existingUser != nil {
		log.Fatal("このログインIDは既に使用されています")
	}

	// パスワードハッシュ化
	hashedPassword, err := utils.HashPassword(*password)
	if err != nil {
		log.Fatal("パスワードハッシュ化エラー:", err)
	}

	// 管理者ユーザー作成
	adminUser := &models.AdminUser{
		LoginID:   *loginID,
		Password:  hashedPassword,
		Name:      *name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	err = adminRepo.Create(adminUser)
	if err != nil {
		log.Fatal("管理者ユーザー作成エラー:", err)
	}

	fmt.Printf("✅ 管理者ユーザーを作成しました\n")
	fmt.Printf("ログインID: %s\n", *loginID)
	fmt.Printf("表示名: %s\n", *name)
	fmt.Printf("ID: %d\n", adminUser.ID)
	fmt.Println("管理ツールにログインできます: http://localhost:8080/login")
}
package main

import (
	"fmt"
	"log"

	"strikepad-manage-tool/config"
)

func main() {
	fmt.Println("=== admin_users テーブルクリーンアップ ===")

	// データベース接続
	dbConfig := config.NewDatabaseConfig()
	db, err := config.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	// admin_usersテーブルをクリア
	result := db.Exec("DELETE FROM admin_users WHERE login_id = 'admin'")
	if result.Error != nil {
		log.Fatal("管理者ユーザー削除エラー:", result.Error)
	}

	fmt.Printf("削除されたレコード数: %d\n", result.RowsAffected)

	// マイグレーション状態をリセット（20250818000001まで戻す）
	result = db.Exec("UPDATE atlas_schema_revisions.atlas_schema_revisions SET version = '20250818000001' WHERE version = '20250818000002'")
	if result.Error != nil {
		// マイグレーション状態テーブルが存在しない場合は無視
		fmt.Printf("マイグレーション状態更新エラー（無視可能）: %v\n", result.Error)
	} else {
		fmt.Printf("マイグレーション状態をリセットしました\n")
	}

	fmt.Println("クリーンアップ完了")
}
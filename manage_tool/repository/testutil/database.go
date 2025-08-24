package testutil

import (
	"log"
	"os"
	"testing"

	"strikepad-manage-tool/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB はテスト用のインメモリSQLiteデータベースを設定
func SetupTestDB(t *testing.T) *gorm.DB {
	// インメモリSQLiteデータベースを使用
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Silent, // テスト時はログを出力しない
			},
		),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// テーブル作成
	err = db.AutoMigrate(
		&models.Maker{},
		&models.Ball{},
		&models.Cover{},
		&models.Core{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// CleanupTestDB はテスト用データベースをクリーンアップ
func CleanupTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// SeedMakersData はテスト用のメーカーデータを挿入
func SeedMakersData(db *gorm.DB) []*models.Maker {
	makers := []*models.Maker{
		{ID: 1, Name: "Storm", IsDeleted: false},
		{ID: 2, Name: "Hammer", IsDeleted: false},
		{ID: 3, Name: "Brunswick", IsDeleted: false},
		{ID: 4, Name: "Deleted Maker", IsDeleted: true}, // 論理削除済み
	}

	for _, maker := range makers {
		db.Create(maker)
	}

	return makers
}

// SeedRelatedData はメーカーに関連するボール、カバー、コアのデータを挿入
func SeedRelatedData(db *gorm.DB) {
	// カバーデータ
	covers := []*models.Cover{
		{ID: 1, Name: "Reactive Resin", MakerID: 1, MaterialType: 1, Rank: 1, IsDeleted: false},
		{ID: 2, Name: "Urethane", MakerID: 1, MaterialType: 2, Rank: 2, IsDeleted: false},
		{ID: 3, Name: "Plastic", MakerID: 2, MaterialType: 3, Rank: 1, IsDeleted: false},
	}

	for _, cover := range covers {
		db.Create(cover)
	}

	// コアデータ
	cores := []*models.Core{
		{ID: 1, Name: "Asymmetric Core", MakerID: 1, RG: 2.5, DeltaRG: 0.05, SymmetryFlag: false, IsDeleted: false},
		{ID: 2, Name: "Symmetric Core", MakerID: 1, RG: 2.4, DeltaRG: 0.03, SymmetryFlag: true, IsDeleted: false},
		{ID: 3, Name: "Light Core", MakerID: 2, RG: 2.6, DeltaRG: 0.04, SymmetryFlag: true, IsDeleted: false},
	}

	for _, core := range cores {
		db.Create(core)
	}

	// ボールデータ
	balls := []*models.Ball{
		{ID: 1, Name: "Storm Ball 1", MakerID: 1, CoreID: 1, CoverID: 1, URL: "http://example.com/1", IsDeleted: false},
		{ID: 2, Name: "Storm Ball 2", MakerID: 1, CoreID: 2, CoverID: 2, URL: "http://example.com/2", IsDeleted: false},
		{ID: 3, Name: "Hammer Ball 1", MakerID: 2, CoreID: 3, CoverID: 3, URL: "http://example.com/3", IsDeleted: false},
	}

	for _, ball := range balls {
		db.Create(ball)
	}
}

// AssertMakerEqual はメーカーオブジェクトを比較
func AssertMakerEqual(t *testing.T, expected, actual *models.Maker) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil || actual == nil {
		t.Errorf("Expected one to be nil, but got expected=%v, actual=%v", expected, actual)
		return
	}

	if expected.ID != actual.ID {
		t.Errorf("Expected ID %d, got %d", expected.ID, actual.ID)
	}
	if expected.Name != actual.Name {
		t.Errorf("Expected Name %s, got %s", expected.Name, actual.Name)
	}
	if expected.LogoFile != actual.LogoFile {
		t.Errorf("Expected LogoFile %s, got %s", expected.LogoFile, actual.LogoFile)
	}
	if expected.IsDeleted != actual.IsDeleted {
		t.Errorf("Expected IsDeleted %v, got %v", expected.IsDeleted, actual.IsDeleted)
	}
}

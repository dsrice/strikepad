package repository

import (
	"database/sql"
	"testing"

	"strikepad-manage-tool/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestHelpers は共通のテストヘルパー関数を提供

// setupMockDB はテスト用のモックデータベースとGORMを設定
// すべてのリポジトリテストで共通利用
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		// ログ出力を無効化してテスト出力をクリーンに
		Logger: nil,
	})
	require.NoError(t, err)

	return db, mock, gormDB
}

// createTestMaker はテスト用のメーカーオブジェクトを作成
func createTestMaker(id uint, name string, logoFile ...string) *models.Maker {
	logo := ""
	if len(logoFile) > 0 {
		logo = logoFile[0]
	}

	return &models.Maker{
		ID:        id,
		Name:      name,
		LogoFile:  logo,
		IsDeleted: false,
	}
}

// createTestCore はテスト用のコアオブジェクトを作成
func createTestCore(id uint, name string, makerID uint, rg float32, deltaRG float32, symmetryFlag bool) *models.Core {
	return &models.Core{
		ID:           id,
		Name:         name,
		MakerID:      makerID,
		RG:           rg,
		DeltaRG:      deltaRG,
		SymmetryFlag: symmetryFlag,
		IsDeleted:    false,
	}
}

// createTestCover はテスト用のカバーオブジェクトを作成
func createTestCover(id uint, name string, makerID uint, materialType int, rank int) *models.Cover {
	return &models.Cover{
		ID:           id,
		Name:         name,
		MakerID:      makerID,
		MaterialType: materialType,
		Rank:         rank,
		IsDeleted:    false,
	}
}

// assertMakerEqual はメーカーオブジェクトの等価性をチェック
func assertMakerEqual(t *testing.T, expected *models.Maker, actual *models.Maker) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID, "ID should match")
	require.Equal(t, expected.Name, actual.Name, "Name should match")
	require.Equal(t, expected.LogoFile, actual.LogoFile, "LogoFile should match")
	require.Equal(t, expected.IsDeleted, actual.IsDeleted, "IsDeleted should match")
}

// assertCoreEqual はコアオブジェクトの等価性をチェック
func assertCoreEqual(t *testing.T, expected *models.Core, actual *models.Core) {
	t.Helper()

	require.Equal(t, expected.ID, actual.ID, "ID should match")
	require.Equal(t, expected.Name, actual.Name, "Name should match")
	require.Equal(t, expected.MakerID, actual.MakerID, "MakerID should match")
	require.Equal(t, expected.RG, actual.RG, "RG should match")
	require.Equal(t, expected.DeltaRG, actual.DeltaRG, "DeltaRG should match")
	require.Equal(t, expected.SymmetryFlag, actual.SymmetryFlag, "SymmetryFlag should match")
	require.Equal(t, expected.IsDeleted, actual.IsDeleted, "IsDeleted should match")

	if expected.InitDiff != nil && actual.InitDiff != nil {
		require.Equal(t, *expected.InitDiff, *actual.InitDiff, "InitDiff should match")
	} else {
		require.Equal(t, expected.InitDiff, actual.InitDiff, "InitDiff pointer should match")
	}
}

// mockRowsForMaker はメーカー用のモック行データを作成
func mockRowsForMaker(makers ...*models.Maker) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "name", "logo_file", "created_at", "updated_at", "is_deleted",
	})

	for _, maker := range makers {
		rows.AddRow(
			maker.ID,
			maker.Name,
			maker.LogoFile,
			maker.CreatedAt,
			maker.UpdatedAt,
			maker.IsDeleted,
		)
	}

	return rows
}

// mockRowsForCore はコア用のモック行データを作成
func mockRowsForCore(cores ...*models.Core) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "name", "rg", "delta_rg", "init_diff", "symmetry_flag",
		"maker_id", "created_at", "updated_at", "is_deleted",
	})

	for _, core := range cores {
		rows.AddRow(
			core.ID,
			core.Name,
			core.RG,
			core.DeltaRG,
			core.InitDiff,
			core.SymmetryFlag,
			core.MakerID,
			core.CreatedAt,
			core.UpdatedAt,
			core.IsDeleted,
		)
	}

	return rows
}

// mockRowsForMakerList はMakerListItem用のモック行データを作成
func mockRowsForMakerList() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "name", "created_at", "updated_at", "ball_count", "cover_count", "core_count",
	})
}

// TestCommonHelperFunctions は共通ヘルパー関数のテスト
func TestCreateTestMaker(t *testing.T) {
	// ロゴファイルなしの場合
	maker1 := createTestMaker(1, "Test Maker")
	require.Equal(t, uint(1), maker1.ID)
	require.Equal(t, "Test Maker", maker1.Name)
	require.Equal(t, "", maker1.LogoFile)
	require.False(t, maker1.IsDeleted)

	// ロゴファイルありの場合
	maker2 := createTestMaker(2, "Test Maker 2", "logo.png")
	require.Equal(t, uint(2), maker2.ID)
	require.Equal(t, "Test Maker 2", maker2.Name)
	require.Equal(t, "logo.png", maker2.LogoFile)
	require.False(t, maker2.IsDeleted)
}

func TestCreateTestCore(t *testing.T) {
	core := createTestCore(1, "Test Core", 2, 2.5, 0.05, true)
	require.Equal(t, uint(1), core.ID)
	require.Equal(t, "Test Core", core.Name)
	require.Equal(t, uint(2), core.MakerID)
	require.Equal(t, float32(2.5), core.RG)
	require.Equal(t, float32(0.05), core.DeltaRG)
	require.True(t, core.SymmetryFlag)
	require.False(t, core.IsDeleted)
}

func TestAssertMakerEqual(t *testing.T) {
	maker1 := createTestMaker(1, "Test", "logo.png")
	maker2 := createTestMaker(1, "Test", "logo.png")

	// 同じ値なので正常終了するはず
	assertMakerEqual(t, maker1, maker2)
}

func TestAssertCoreEqual(t *testing.T) {
	core1 := createTestCore(1, "Test Core", 2, 2.5, 0.05, true)
	core2 := createTestCore(1, "Test Core", 2, 2.5, 0.05, true)

	// 同じ値なので正常終了するはず
	assertCoreEqual(t, core1, core2)
}

func TestMockRowsForMaker(t *testing.T) {
	maker1 := createTestMaker(1, "Maker 1")
	maker2 := createTestMaker(2, "Maker 2", "logo.png")

	rows := mockRowsForMaker(maker1, maker2)
	require.NotNil(t, rows)

	// mockRowsは作成できることをテスト
	// 実際のデータチェックはsqlmockが実行時に行う
}

func TestMockRowsForCore(t *testing.T) {
	core1 := createTestCore(1, "Core 1", 1, 2.5, 0.05, true)
	core2 := createTestCore(2, "Core 2", 2, 2.6, 0.04, false)

	rows := mockRowsForCore(core1, core2)
	require.NotNil(t, rows)
}
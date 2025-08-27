package repository_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCoreRepository_Create(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	core := &models.Core{
		Name:         "Test Core",
		MakerID:      1,
		RG:           2.5,
		DeltaRG:      0.05,
		SymmetryFlag: true,
	}

	// モックのInsert期待値を設定
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "cores"`)).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // init_diff
			sqlmock.AnyArg(), // deleted_at
			"Test Core",      // name
			uint(1),          // maker_id
			float32(2.5),     // rg
			float32(0.05),    // delta_rg
			true,             // symmetry_flag
			false,            // is_deleted
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(core)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), core.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_GetByID(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	now := time.Now()
	expectedCore := &models.Core{
		ID:           1,
		Name:         "Test Core",
		RG:           2.5,
		DeltaRG:      0.05,
		SymmetryFlag: true,
		MakerID:      1,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsDeleted:    false,
	}

	// コア取得用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cores"`)).
		WithArgs(uint(1), false, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "rg", "delta_rg", "init_diff", "symmetry_flag",
			"maker_id", "created_at", "updated_at", "is_deleted",
		}).AddRow(
			1, "Test Core", 2.5, 0.05, nil, true, 1, now, now, false,
		))

	// メーカー情報のPreload用モック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "makers"`)).
		WithArgs(uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "created_at", "updated_at", "is_deleted",
		}).AddRow(
			1, "Test Maker", now, now, false,
		))

	result, err := repo.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCore.ID, result.ID)
	assert.Equal(t, expectedCore.Name, result.Name)
	assert.Equal(t, expectedCore.RG, result.RG)
	assert.Equal(t, expectedCore.DeltaRG, result.DeltaRG)
	assert.Equal(t, expectedCore.SymmetryFlag, result.SymmetryFlag)
	assert.Equal(t, expectedCore.MakerID, result.MakerID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_GetByID_NotFound(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	// レコードが見つからない場合のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cores"`)).
		WithArgs(uint(999), false, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(999)

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_GetAll(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	now := time.Now()

	// GetAll用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cores"`)).
		WithArgs(false, 10).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "rg", "delta_rg", "init_diff", "symmetry_flag",
			"maker_id", "created_at", "updated_at", "is_deleted",
		}).
			AddRow(1, "Core 1", 2.5, 0.05, nil, true, 1, now, now, false).
			AddRow(2, "Core 2", 2.6, 0.04, nil, false, 2, now, now, false))

	// Preload用のメーカー情報モック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "makers"`)).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "created_at", "updated_at", "is_deleted",
		}).
			AddRow(1, "Maker 1", now, now, false).
			AddRow(2, "Maker 2", now, now, false))

	cores, err := repo.GetAll(0, 10)

	assert.NoError(t, err)
	assert.Len(t, cores, 2)
	assert.Equal(t, "Core 1", cores[0].Name)
	assert.Equal(t, "Core 2", cores[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_GetByMakerID(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	now := time.Now()

	// GetByMakerID用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cores"`)).
		WithArgs(uint(1), false, 5).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "rg", "delta_rg", "init_diff", "symmetry_flag",
			"maker_id", "created_at", "updated_at", "is_deleted",
		}).
			AddRow(1, "Core 1", 2.5, 0.05, nil, true, 1, now, now, false))

	// Preload用のメーカー情報モック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "makers"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "created_at", "updated_at", "is_deleted",
		}).
			AddRow(1, "Maker 1", now, now, false))

	cores, err := repo.GetByMakerID(1, 0, 5)

	assert.NoError(t, err)
	assert.Len(t, cores, 1)
	assert.Equal(t, "Core 1", cores[0].Name)
	assert.Equal(t, uint(1), cores[0].MakerID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_Search(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	now := time.Now()

	// Search用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cores"`)).
		WithArgs("%test%", false, 10).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "rg", "delta_rg", "init_diff", "symmetry_flag",
			"maker_id", "created_at", "updated_at", "is_deleted",
		}).
			AddRow(1, "Test Core", 2.5, 0.05, nil, true, 1, now, now, false))

	// Preload用のメーカー情報モック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "makers"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "created_at", "updated_at", "is_deleted",
		}).
			AddRow(1, "Maker 1", now, now, false))

	cores, err := repo.Search("test", 0, 10)

	assert.NoError(t, err)
	assert.Len(t, cores, 1)
	assert.Equal(t, "Test Core", cores[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_Count(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	// Count用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "cores"`)).
		WithArgs(false).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.Count()

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_CountByMaker(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	// CountByMaker用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "cores"`)).
		WithArgs(uint(1), false).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	count, err := repo.CountByMaker(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_Update(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	core := &models.Core{
		ID:           1,
		Name:         "Updated Core",
		MakerID:      2,
		RG:           2.7,
		DeltaRG:      0.06,
		SymmetryFlag: false,
	}

	// Update用のモック
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "cores" SET`)).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // init_diff
			sqlmock.AnyArg(), // deleted_at
			"Updated Core",   // name
			uint(2),          // maker_id
			float32(2.7),     // rg
			float32(0.06),    // delta_rg
			false,            // symmetry_flag
			false,            // is_deleted
			uint(1),          // id (WHERE condition)
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(core)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_Delete(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	// Delete用のモック（論理削除）
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "cores" SET`)).
		WithArgs(sqlmock.AnyArg(), true, sqlmock.AnyArg(), uint(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_Create_Error(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	core := &models.Core{
		Name:         "Test Core",
		MakerID:      1,
		RG:           2.5,
		DeltaRG:      0.05,
		SymmetryFlag: true,
	}

	// エラーが発生する場合のモック
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "cores"`)).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	err := repo.Create(core)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCoreRepository_GetByID_Error(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewCoreRepositoryInterface(gormDB)

	// データベースエラーの場合のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cores"`)).
		WithArgs(uint(1), false, 1).
		WillReturnError(errors.New("database connection error"))

	result, err := repo.GetByID(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
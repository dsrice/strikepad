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

func TestMakerRepository_GetAllSimple(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	now := time.Now()

	// GetAllSimple用のモック（GORMが生成する実際のクエリに合わせる）
	mock.ExpectQuery(`SELECT \* FROM "makers"`).
		WithArgs(false).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "logo_file", "created_at", "updated_at", "is_deleted", "deleted_at",
		}).
			AddRow(1, "Brunswick", "brunswick.png", now, now, false, nil).
			AddRow(2, "Hammer", "hammer.png", now, now, false, nil).
			AddRow(3, "Storm", "storm.png", now, now, false, nil))

	makers, err := repo.GetAllSimple()

	assert.NoError(t, err)
	assert.Len(t, makers, 3)
	assert.Equal(t, "Brunswick", makers[0].Name) // ORDER BY name ASCで最初
	assert.Equal(t, "Hammer", makers[1].Name)
	assert.Equal(t, "Storm", makers[2].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_GetAllSimple_Empty(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	// 空の結果を返すモック
	mock.ExpectQuery(`SELECT \* FROM "makers"`).
		WithArgs(false).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "logo_file", "created_at", "updated_at", "is_deleted", "deleted_at",
		}))

	makers, err := repo.GetAllSimple()

	assert.NoError(t, err)
	assert.Len(t, makers, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_Create_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	maker := &models.Maker{
		Name:     "Test Maker",
		LogoFile: "test_logo.png",
	}

	// Create用のモック
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "makers"`)).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			"Test Maker",     // name
			"test_logo.png",  // logo_file
			false,            // is_deleted
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(maker)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), maker.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_GetByID_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	now := time.Now()

	// GetByID用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "makers"`)).
		WithArgs(uint(1), false).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "logo_file", "created_at", "updated_at", "is_deleted",
		}).AddRow(
			1, "Test Maker", "test_logo.png", now, now, false,
		))

	maker, err := repo.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, maker)
	assert.Equal(t, uint(1), maker.ID)
	assert.Equal(t, "Test Maker", maker.Name)
	assert.Equal(t, "test_logo.png", maker.LogoFile)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_GetByID_NotFound_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	// レコードが見つからない場合のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "makers"`)).
		WithArgs(uint(999), false).
		WillReturnError(gorm.ErrRecordNotFound)

	maker, err := repo.GetByID(999)

	assert.NoError(t, err)
	assert.Nil(t, maker)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_GetByName_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	now := time.Now()

	// GetByName用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "makers"`)).
		WithArgs("Test Maker", false).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "logo_file", "created_at", "updated_at", "is_deleted",
		}).AddRow(
			1, "Test Maker", "test_logo.png", now, now, false,
		))

	maker, err := repo.GetByName("Test Maker")

	assert.NoError(t, err)
	assert.NotNil(t, maker)
	assert.Equal(t, "Test Maker", maker.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_Update_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	maker := &models.Maker{
		ID:       1,
		Name:     "Updated Maker",
		LogoFile: "updated_logo.png",
	}

	// Update用のモック
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "makers"`)).
		WithArgs(
			sqlmock.AnyArg(),   // updated_at
			sqlmock.AnyArg(),   // deleted_at
			"Updated Maker",    // name
			"updated_logo.png", // logo_file
			false,              // is_deleted
			uint(1),            // id (WHERE condition)
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Update(maker)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_UpdateLogoFile_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	// UpdateLogoFile用のモック
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "makers"`)).
		WithArgs("new_logo.png", uint(1), false).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateLogoFile(1, "new_logo.png")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_Delete_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	// Delete用のモック（論理削除）
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "makers"`)).
		WithArgs(true, sqlmock.AnyArg(), uint(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_GetTotalMakersCount_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	// GetTotalMakersCount用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "makers"`)).
		WithArgs(false).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.GetTotalMakersCount()

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_GetMakerStats_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	// TotalMakers用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "makers"`)).
		WithArgs(false).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	// ActiveMakers用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(DISTINCT "makers"."id") FROM "makers"`)).
		WithArgs(false, false).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// TotalBalls用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "balls"`)).
		WithArgs(false).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	// TotalCovers用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "covers"`)).
		WithArgs(false).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4))

	// TotalCores用のモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "cores"`)).
		WithArgs(false).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(6))

	stats, err := repo.GetMakerStats()

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(3), stats.TotalMakers)
	assert.Equal(t, int64(2), stats.ActiveMakers)
	assert.Equal(t, int64(5), stats.TotalBalls)
	assert.Equal(t, int64(4), stats.TotalCovers)
	assert.Equal(t, int64(6), stats.TotalCores)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_SearchMakers_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	now := time.Now()

	// SearchMakers用のRawクエリモック
	expectedQuery := `
		SELECT 
			m.id,
			m.name,
			m.created_at,
			m.updated_at,
			COUNT(DISTINCT b.id) as ball_count,
			COUNT(DISTINCT c.id) as cover_count,
			COUNT(DISTINCT cr.id) as core_count
		FROM makers m
		LEFT JOIN balls b ON m.id = b.maker_id AND b.is_deleted = false
		LEFT JOIN covers c ON m.id = c.maker_id AND c.is_deleted = false
		LEFT JOIN cores cr ON m.id = cr.maker_id AND cr.is_deleted = false
		WHERE m.is_deleted = false AND m.name ILIKE ?
		GROUP BY m.id, m.name, m.created_at, m.updated_at
		ORDER BY m.id ASC
		LIMIT ? OFFSET ?`

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs("%storm%", 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "created_at", "updated_at", "ball_count", "cover_count", "core_count",
		}).AddRow(
			1, "Storm", now, now, 2, 1, 1,
		))

	makers, err := repo.SearchMakers("storm", 0, 10)

	assert.NoError(t, err)
	assert.Len(t, makers, 1)
	assert.Equal(t, "Storm", makers[0].Name)
	assert.Equal(t, int64(2), makers[0].BallCount)
	assert.Equal(t, int64(1), makers[0].CoverCount)
	assert.Equal(t, int64(1), makers[0].CoreCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_GetAllMakers_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	now := time.Now()

	// GetAllMakers用のRawクエリモック
	expectedQuery := `
		SELECT 
			m.id,
			m.name,
			m.created_at,
			m.updated_at,
			COUNT(DISTINCT b.id) as ball_count,
			COUNT(DISTINCT c.id) as cover_count,
			COUNT(DISTINCT cr.id) as core_count
		FROM makers m
		LEFT JOIN balls b ON m.id = b.maker_id AND b.is_deleted = false
		LEFT JOIN covers c ON m.id = c.maker_id AND c.is_deleted = false
		LEFT JOIN cores cr ON m.id = cr.maker_id AND cr.is_deleted = false
		WHERE m.is_deleted = false
		GROUP BY m.id, m.name, m.created_at, m.updated_at
		ORDER BY m.id ASC
		LIMIT ? OFFSET ?`

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "created_at", "updated_at", "ball_count", "cover_count", "core_count",
		}).
			AddRow(1, "Brunswick", now, now, 3, 2, 1).
			AddRow(2, "Storm", now, now, 2, 1, 2))

	makers, err := repo.GetAllMakers(0, 10)

	assert.NoError(t, err)
	assert.Len(t, makers, 2)
	assert.Equal(t, "Brunswick", makers[0].Name)
	assert.Equal(t, int64(3), makers[0].BallCount)
	assert.Equal(t, "Storm", makers[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// エラーケースのテスト
func TestMakerRepository_GetAllSimple_Error(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	// データベースエラーを返すモック
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "makers"`)).
		WithArgs(false).
		WillReturnError(errors.New("database connection error"))

	makers, err := repo.GetAllSimple()

	assert.Error(t, err)
	assert.Nil(t, makers)
	assert.Contains(t, err.Error(), "database connection error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMakerRepository_Create_Error_SqlMock(t *testing.T) {
	db, mock, gormDB := setupMockDB(t)
	defer db.Close()

	repo := repository.NewMakerRepository(gormDB)

	maker := &models.Maker{
		Name:     "Test Maker",
		LogoFile: "test_logo.png",
	}

	// エラーが発生する場合のモック
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "makers"`)).
		WillReturnError(errors.New("unique constraint violation"))
	mock.ExpectRollback()

	err := repo.Create(maker)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unique constraint violation")
	assert.NoError(t, mock.ExpectationsWereMet())
}
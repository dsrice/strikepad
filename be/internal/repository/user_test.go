package repository

import (
	"testing"
	"time"

	"strikepad-backend/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const testEmail = "test@example.com"

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.db = gormDB
	suite.mock = mock
	suite.repo = NewUserRepository(gormDB)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	err := suite.mock.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	email := testEmail
	user := &model.User{
		ProviderType: "email",
		DisplayName:  "Test User",
		Email:        &email,
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("INSERT INTO `users`").
		WithArgs("email", nil, "test@example.com", "Test User", nil, false, false, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	createdUser, err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "email", createdUser.ProviderType)
	assert.Equal(suite.T(), "Test User", createdUser.DisplayName)
}

func (suite *UserRepositoryTestSuite) TestGetByID() {
	email := testEmail
	now := time.Now()

	suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
			AddRow(1, "email", nil, email, "Test User", nil, false, now, now, false, nil))

	found, err := suite.repo.GetByID(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test User", found.DisplayName)
	assert.Equal(suite.T(), &email, found.Email)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	email := testEmail
	now := time.Now()

	suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(email, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
			AddRow(1, "email", nil, email, "Test User", nil, false, now, now, false, nil))

	found, err := suite.repo.GetByEmail(email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test User", found.DisplayName)
	assert.Equal(suite.T(), uint(1), found.ID)
}

func (suite *UserRepositoryTestSuite) TestList() {
	email1 := "user1@example.com"
	email2 := "user2@example.com"
	now := time.Now()

	suite.mock.ExpectQuery("SELECT \\* FROM `users`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
			AddRow(1, "email", nil, email1, "User 1", nil, false, now, now, false, nil).
			AddRow(2, "email", nil, email2, "User 2", nil, false, now, now, false, nil))

	result, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
}

func (suite *UserRepositoryTestSuite) TestDelete() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("DELETE FROM `users` WHERE `users`.`id` = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Delete(1)
	assert.NoError(suite.T(), err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

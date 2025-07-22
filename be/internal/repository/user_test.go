package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strikepad-backend/internal/model"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)

	err = db.AutoMigrate(&model.User{})
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.repo = NewUserRepository(db)
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), user.ID)
}

func (suite *UserRepositoryTestSuite) TestGetByID() {
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	found, err := suite.repo.GetByID(user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Name, found.Name)
	assert.Equal(suite.T(), user.Email, found.Email)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	found, err := suite.repo.GetByEmail(user.Email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Name, found.Name)
	assert.Equal(suite.T(), user.ID, found.ID)
}

func (suite *UserRepositoryTestSuite) TestList() {
	users := []*model.User{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
	}

	for _, user := range users {
		err := suite.repo.Create(user)
		assert.NoError(suite.T(), err)
	}

	result, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
}

func (suite *UserRepositoryTestSuite) TestDelete() {
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	err := suite.repo.Create(user)
	assert.NoError(suite.T(), err)

	err = suite.repo.Delete(user.ID)
	assert.NoError(suite.T(), err)

	_, err = suite.repo.GetByID(user.ID)
	assert.Error(suite.T(), err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
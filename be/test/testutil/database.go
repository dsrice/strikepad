package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strikepad-backend/internal/model"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.User{})
	assert.NoError(t, err)

	return db
}

func CleanupTestDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func CreateTestUser(t *testing.T, db *gorm.DB, name, email string) *model.User {
	user := &model.User{
		Name:  name,
		Email: email,
	}

	err := db.Create(user).Error
	assert.NoError(t, err)

	return user
}
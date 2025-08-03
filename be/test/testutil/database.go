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
	sqlDB, err := db.DB()
	if err == nil && sqlDB != nil {
		_ = sqlDB.Close()
	}
}

func CreateTestUser(t *testing.T, db *gorm.DB, displayName, email string) *model.User {
	user := &model.User{
		ProviderType: "email",
		DisplayName:  displayName,
		Email:        &email,
	}

	err := db.Create(user).Error
	assert.NoError(t, err)

	return user
}

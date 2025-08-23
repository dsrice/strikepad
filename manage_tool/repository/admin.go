package repository

import (
	"errors"

	"strikepad-manage-tool/models"

	"gorm.io/gorm"
)

// AdminRepository は管理者ユーザー関連のデータベース操作を行う
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository は新しい管理者リポジトリを作成
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetByLoginID はログインIDで管理者ユーザーを取得
func (r *AdminRepository) GetByLoginID(loginID string) (*models.AdminUser, error) {
	var adminUser models.AdminUser

	result := r.db.Where("login_id = ? AND is_deleted = ?", loginID, false).First(&adminUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // ユーザーが見つからない場合はnilを返す
		}
		return nil, result.Error
	}

	return &adminUser, nil
}

// GetByID はIDで管理者ユーザーを取得
func (r *AdminRepository) GetByID(id uint) (*models.AdminUser, error) {
	var adminUser models.AdminUser

	result := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&adminUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &adminUser, nil
}

// Create は新しい管理者ユーザーを作成
func (r *AdminRepository) Create(adminUser *models.AdminUser) error {
	return r.db.Create(adminUser).Error
}

// Update は管理者ユーザー情報を更新
func (r *AdminRepository) Update(adminUser *models.AdminUser) error {
	return r.db.Save(adminUser).Error
}

// Delete は管理者ユーザーを論理削除
func (r *AdminRepository) Delete(id uint) error {
	return r.db.Model(&models.AdminUser{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": gorm.DeletedAt{},
	}).Error
}

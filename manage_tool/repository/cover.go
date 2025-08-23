package repository

import (
	"errors"

	"strikepad-manage-tool/models"

	"gorm.io/gorm"
)

// CoverRepository はカバー関連のデータベース操作を行う
type CoverRepository struct {
	db *gorm.DB
}

// NewCoverRepository は新しいカバーリポジトリを作成
func NewCoverRepository(db *gorm.DB) *CoverRepository {
	return &CoverRepository{db: db}
}

// GetAll は全てのカバーを取得
func (r *CoverRepository) GetAll(offset, limit int) ([]*models.Cover, error) {
	var covers []*models.Cover

	query := r.db.Where("is_deleted = ?", false).
		Preload("Maker").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&covers)
	if result.Error != nil {
		return nil, result.Error
	}

	return covers, nil
}

// GetByID はIDでカバーを取得
func (r *CoverRepository) GetByID(id uint) (*models.Cover, error) {
	var cover models.Cover

	result := r.db.Where("id = ? AND is_deleted = ?", id, false).
		Preload("Maker").
		First(&cover)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &cover, nil
}

// GetByMakerID はメーカーIDでカバーを取得
func (r *CoverRepository) GetByMakerID(makerID uint, offset, limit int) ([]*models.Cover, error) {
	var covers []*models.Cover

	query := r.db.Where("maker_id = ? AND is_deleted = ?", makerID, false).
		Preload("Maker").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&covers)
	if result.Error != nil {
		return nil, result.Error
	}

	return covers, nil
}

// Search は名前でカバーを検索
func (r *CoverRepository) Search(keyword string, offset, limit int) ([]*models.Cover, error) {
	var covers []*models.Cover

	query := r.db.Where("name ILIKE ? AND is_deleted = ?", "%"+keyword+"%", false).
		Preload("Maker").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&covers)
	if result.Error != nil {
		return nil, result.Error
	}

	return covers, nil
}

// Count は総カバー数を取得
func (r *CoverRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&models.Cover{}).Where("is_deleted = ?", false).Count(&count)
	return count, result.Error
}

// CountByMaker はメーカー別のカバー数を取得
func (r *CoverRepository) CountByMaker(makerID uint) (int64, error) {
	var count int64
	result := r.db.Model(&models.Cover{}).
		Where("maker_id = ? AND is_deleted = ?", makerID, false).
		Count(&count)
	return count, result.Error
}

// Create は新しいカバーを作成
func (r *CoverRepository) Create(cover *models.Cover) error {
	return r.db.Create(cover).Error
}

// Update はカバー情報を更新
func (r *CoverRepository) Update(cover *models.Cover) error {
	return r.db.Save(cover).Error
}

// Delete はカバーを論理削除
func (r *CoverRepository) Delete(id uint) error {
	return r.db.Model(&models.Cover{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.DeletedAt{},
		}).Error
}

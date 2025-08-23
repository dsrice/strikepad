package repository

import (
	"errors"

	"strikepad-manage-tool/models"

	"gorm.io/gorm"
)

// CoreRepository はコア関連のデータベース操作を行う
type CoreRepository struct {
	db *gorm.DB
}

// NewCoreRepository は新しいコアリポジトリを作成
func NewCoreRepository(db *gorm.DB) *CoreRepository {
	return &CoreRepository{db: db}
}

// GetAll は全てのコアを取得
func (r *CoreRepository) GetAll(offset, limit int) ([]*models.Core, error) {
	var cores []*models.Core

	query := r.db.Where("is_deleted = ?", false).
		Preload("Maker").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&cores)
	if result.Error != nil {
		return nil, result.Error
	}

	return cores, nil
}

// GetByID はIDでコアを取得
func (r *CoreRepository) GetByID(id uint) (*models.Core, error) {
	var core models.Core

	result := r.db.Where("id = ? AND is_deleted = ?", id, false).
		Preload("Maker").
		First(&core)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &core, nil
}

// GetByMakerID はメーカーIDでコアを取得
func (r *CoreRepository) GetByMakerID(makerID uint, offset, limit int) ([]*models.Core, error) {
	var cores []*models.Core

	query := r.db.Where("maker_id = ? AND is_deleted = ?", makerID, false).
		Preload("Maker").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&cores)
	if result.Error != nil {
		return nil, result.Error
	}

	return cores, nil
}

// Search は名前でコアを検索
func (r *CoreRepository) Search(keyword string, offset, limit int) ([]*models.Core, error) {
	var cores []*models.Core

	query := r.db.Where("name ILIKE ? AND is_deleted = ?", "%"+keyword+"%", false).
		Preload("Maker").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&cores)
	if result.Error != nil {
		return nil, result.Error
	}

	return cores, nil
}

// Count は総コア数を取得
func (r *CoreRepository) Count() (int64, error) {
	var count int64
	result := r.db.Model(&models.Core{}).Where("is_deleted = ?", false).Count(&count)
	return count, result.Error
}

// CountByMaker はメーカー別のコア数を取得
func (r *CoreRepository) CountByMaker(makerID uint) (int64, error) {
	var count int64
	result := r.db.Model(&models.Core{}).
		Where("maker_id = ? AND is_deleted = ?", makerID, false).
		Count(&count)
	return count, result.Error
}

// Create は新しいコアを作成
func (r *CoreRepository) Create(core *models.Core) error {
	return r.db.Create(core).Error
}

// Update はコア情報を更新
func (r *CoreRepository) Update(core *models.Core) error {
	return r.db.Save(core).Error
}

// Delete はコアを論理削除
func (r *CoreRepository) Delete(id uint) error {
	return r.db.Model(&models.Core{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.DeletedAt{},
		}).Error
}

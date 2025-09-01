package repository

import (
	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository/ri"

	"gorm.io/gorm"
)

// ballRepository はボールリポジトリの実装
type ballRepository struct {
	db *gorm.DB
}

// NewBallRepositoryInterface はDI用のBallRepositoryInterfaceを返す
func NewBallRepositoryInterface(db *gorm.DB) ri.BallRepositoryInterface {
	return &ballRepository{db: db}
}

// GetAll は全てのボールを取得（ページネーション付き）
func (r *ballRepository) GetAll(offset, limit int) ([]*models.Ball, error) {
	var balls []*models.Ball
	err := r.db.Preload("Maker").
		Preload("Core").
		Preload("Core.Maker").
		Preload("Cover").
		Preload("Cover.Maker").
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&balls).Error
	return balls, err
}

// GetByID はIDでボールを取得
func (r *ballRepository) GetByID(id uint) (*models.Ball, error) {
	var ball models.Ball
	err := r.db.Preload("Maker").
		Preload("Core").
		Preload("Core.Maker").
		Preload("Cover").
		Preload("Cover.Maker").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&ball).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ball, err
}

// Create は新しいボールを作成
func (r *ballRepository) Create(ball *models.Ball) error {
	return r.db.Create(ball).Error
}

// Update はボール情報を更新
func (r *ballRepository) Update(ball *models.Ball) error {
	return r.db.Save(ball).Error
}

// Delete はボールを論理削除
func (r *ballRepository) Delete(id uint) error {
	return r.db.Model(&models.Ball{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
		}).Error
}

// Count は全ボール数を取得
func (r *ballRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Ball{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}

// Search はボール名で検索
func (r *ballRepository) Search(query string, offset, limit int) ([]*models.Ball, error) {
	var balls []*models.Ball
	err := r.db.Preload("Maker").
		Preload("Core").
		Preload("Core.Maker").
		Preload("Cover").
		Preload("Cover.Maker").
		Where("is_deleted = ? AND name ILIKE ?", false, "%"+query+"%").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&balls).Error
	return balls, err
}

// GetByMakerID はメーカーIDでボールを取得
func (r *ballRepository) GetByMakerID(makerID uint, offset, limit int) ([]*models.Ball, error) {
	var balls []*models.Ball
	err := r.db.Preload("Maker").
		Preload("Core").
		Preload("Core.Maker").
		Preload("Cover").
		Preload("Cover.Maker").
		Where("is_deleted = ? AND maker_id = ?", false, makerID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&balls).Error
	return balls, err
}

// CountByMaker はメーカー別のボール数を取得
func (r *ballRepository) CountByMaker(makerID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Ball{}).
		Where("is_deleted = ? AND maker_id = ?", false, makerID).
		Count(&count).Error
	return count, err
}
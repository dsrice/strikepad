package repository

import (
	"errors"

	"strikepad-manage-tool/models"

	"gorm.io/gorm"
)

// MakerRepository はメーカー関連のデータベース操作を行う
type MakerRepository struct {
	db *gorm.DB
}

// NewMakerRepository は新しいメーカーリポジトリを作成
func NewMakerRepository(db *gorm.DB) *MakerRepository {
	return &MakerRepository{db: db}
}

// GetAllMakers は全メーカーを取得（ページネーション対応）
func (r *MakerRepository) GetAllMakers(offset, limit int) ([]models.MakerListItem, error) {
	var makers []models.MakerListItem

	query := `
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
		LIMIT ? OFFSET ?
	`

	err := r.db.Raw(query, limit, offset).Scan(&makers).Error
	return makers, err
}

// SearchMakers はメーカーを検索
func (r *MakerRepository) SearchMakers(searchQuery string, offset, limit int) ([]models.MakerListItem, error) {
	var makers []models.MakerListItem

	query := `
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
		LIMIT ? OFFSET ?
	`

	searchPattern := "%" + searchQuery + "%"
	err := r.db.Raw(query, searchPattern, limit, offset).Scan(&makers).Error
	return makers, err
}

// GetByID はIDでメーカーを取得
func (r *MakerRepository) GetByID(id uint) (*models.Maker, error) {
	var maker models.Maker

	result := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&maker)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &maker, nil
}

// GetByName は名前でメーカーを取得
func (r *MakerRepository) GetByName(name string) (*models.Maker, error) {
	var maker models.Maker

	result := r.db.Where("name = ? AND is_deleted = ?", name, false).First(&maker)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &maker, nil
}

// Create は新しいメーカーを作成
func (r *MakerRepository) Create(maker *models.Maker) error {
	return r.db.Create(maker).Error
}

// Update はメーカー情報を更新
func (r *MakerRepository) Update(maker *models.Maker) error {
	return r.db.Save(maker).Error
}

// UpdateLogoFile はメーカーのロゴファイル名を更新
func (r *MakerRepository) UpdateLogoFile(makerID uint, filename string) error {
	return r.db.Model(&models.Maker{}).
		Where("id = ? AND is_deleted = ?", makerID, false).
		Update("logo_file", filename).Error
}

// Delete はメーカーを論理削除
func (r *MakerRepository) Delete(id uint) error {
	return r.db.Model(&models.Maker{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": gorm.DeletedAt{},
	}).Error
}

// GetTotalMakersCount は総メーカー数を取得
func (r *MakerRepository) GetTotalMakersCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Maker{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}

// GetMakerStats はメーカー統計を取得
func (r *MakerRepository) GetMakerStats() (*models.MakerStats, error) {
	stats := &models.MakerStats{}

	// 総メーカー数
	err := r.db.Model(&models.Maker{}).Where("is_deleted = ?", false).Count(&stats.TotalMakers).Error
	if err != nil {
		return nil, err
	}

	// アクティブメーカー数（ボールを持つメーカー）
	err = r.db.Model(&models.Maker{}).
		Joins("INNER JOIN balls ON makers.id = balls.maker_id AND balls.is_deleted = false").
		Where("makers.is_deleted = ?", false).
		Distinct("makers.id").
		Count(&stats.ActiveMakers).Error
	if err != nil {
		return nil, err
	}

	// 総ボール数
	err = r.db.Model(&models.Ball{}).Where("is_deleted = ?", false).Count(&stats.TotalBalls).Error
	if err != nil {
		return nil, err
	}

	// 総カバー数
	err = r.db.Model(&models.Cover{}).Where("is_deleted = ?", false).Count(&stats.TotalCovers).Error
	if err != nil {
		return nil, err
	}

	// 総コア数
	err = r.db.Model(&models.Core{}).Where("is_deleted = ?", false).Count(&stats.TotalCores).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}

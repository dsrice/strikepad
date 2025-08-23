package handlers

import "strikepad-manage-tool/models"

// MakerRepositoryInterface はメーカーリポジトリのインターフェース
type MakerRepositoryInterface interface {
	GetAllMakers(offset, limit int) ([]models.MakerListItem, error)
	SearchMakers(searchQuery string, offset, limit int) ([]models.MakerListItem, error)
	GetByID(id uint) (*models.Maker, error)
	GetByName(name string) (*models.Maker, error)
	Create(maker *models.Maker) error
	Update(maker *models.Maker) error
	UpdateLogoFile(makerID uint, filename string) error
	Delete(id uint) error
	GetTotalMakersCount() (int64, error)
	GetMakerStats() (*models.MakerStats, error)
}
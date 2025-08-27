package ri

import "strikepad-manage-tool/models"

// UserRepositoryInterface はユーザーリポジトリのインターフェース
type UserRepositoryInterface interface {
	GetAllUsers(offset, limit int) ([]models.UserListItem, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUserStatus(id uint, isActive bool) error
	DeleteUser(id uint) error
	GetUserStats() (*models.UserStats, error)
	SearchUsers(query string, offset, limit int) ([]models.UserListItem, error)
	GetTotalUsersCount() (int64, error)
}

// AdminRepositoryInterface は管理者リポジトリのインターフェース
type AdminRepositoryInterface interface {
	GetByLoginID(loginID string) (*models.AdminUser, error)
	GetByID(id uint) (*models.AdminUser, error)
	Create(adminUser *models.AdminUser) error
	Update(adminUser *models.AdminUser) error
	Delete(id uint) error
}

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
	GetAllSimple() ([]*models.Maker, error)
	GetMakerStats() (*models.MakerStats, error)
}

// CoreRepositoryInterface はコアリポジトリのインターフェース
type CoreRepositoryInterface interface {
	GetAll(offset, limit int) ([]*models.Core, error)
	GetByID(id uint) (*models.Core, error)
	GetByMakerID(makerID uint, offset, limit int) ([]*models.Core, error)
	Search(keyword string, offset, limit int) ([]*models.Core, error)
	Count() (int64, error)
	CountByMaker(makerID uint) (int64, error)
	Create(core *models.Core) error
	Update(core *models.Core) error
	Delete(id uint) error
}

// CoverRepositoryInterface はカバーリポジトリのインターフェース
type CoverRepositoryInterface interface {
	GetAll(offset, limit int) ([]*models.Cover, error)
	GetByID(id uint) (*models.Cover, error)
	GetByMakerID(makerID uint, offset, limit int) ([]*models.Cover, error)
	Search(keyword string, offset, limit int) ([]*models.Cover, error)
	Count() (int64, error)
	CountByMaker(makerID uint) (int64, error)
	Create(cover *models.Cover) error
	Update(cover *models.Cover) error
	Delete(id uint) error
}
package mocks

import (
	"strikepad-manage-tool/models"

	"github.com/stretchr/testify/mock"
)

// MockMakerRepository は MakerRepositoryInterface のモック
type MockMakerRepository struct {
	mock.Mock
}

func (m *MockMakerRepository) GetAllMakers(offset, limit int) ([]models.MakerListItem, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.MakerListItem), args.Error(1)
}

func (m *MockMakerRepository) SearchMakers(searchQuery string, offset, limit int) ([]models.MakerListItem, error) {
	args := m.Called(searchQuery, offset, limit)
	return args.Get(0).([]models.MakerListItem), args.Error(1)
}

func (m *MockMakerRepository) GetByID(id uint) (*models.Maker, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Maker), args.Error(1)
}

func (m *MockMakerRepository) GetByName(name string) (*models.Maker, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Maker), args.Error(1)
}

func (m *MockMakerRepository) Create(maker *models.Maker) error {
	args := m.Called(maker)
	return args.Error(0)
}

func (m *MockMakerRepository) Update(maker *models.Maker) error {
	args := m.Called(maker)
	return args.Error(0)
}

func (m *MockMakerRepository) UpdateLogoFile(makerID uint, filename string) error {
	args := m.Called(makerID, filename)
	return args.Error(0)
}

func (m *MockMakerRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMakerRepository) GetTotalMakersCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMakerRepository) GetMakerStats() (*models.MakerStats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MakerStats), args.Error(1)
}

func (m *MockMakerRepository) GetAllSimple() ([]*models.Maker, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Maker), args.Error(1)
}

package repository

import (
	"strikepad-backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List() ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) (*model.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) List() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

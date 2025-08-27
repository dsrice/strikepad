package repository

import (
	"time"

	"strikepad-manage-tool/models"
	"strikepad-manage-tool/repository/ri"

	"gorm.io/gorm"
)

// UserRepository はユーザー操作のリポジトリ
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository は新しいUserRepositoryを作成
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// NewUserRepositoryInterface はDI用のUserRepositoryInterfaceを返す
func NewUserRepositoryInterface(db *gorm.DB) ri.UserRepositoryInterface {
	return NewUserRepository(db)
}

// GetAllUsers は全ユーザーを取得
func (r *UserRepository) GetAllUsers(offset, limit int) ([]models.UserListItem, error) {
	var users []models.UserListItem

	err := r.db.Table("users").
		Select("id, username, email, first_name, last_name, is_active, last_login_at, created_at").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, err
}

// GetUserByID はIDでユーザーを取得
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserStatus はユーザーのアクティブ状態を更新
func (r *UserRepository) UpdateUserStatus(id uint, isActive bool) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

// DeleteUser はユーザーを論理削除
func (r *UserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// GetUserStats はユーザー統計を取得
func (r *UserRepository) GetUserStats() (*models.UserStats, error) {
	stats := &models.UserStats{}

	// 総ユーザー数
	err := r.db.Model(&models.User{}).Count(&stats.TotalUsers).Error
	if err != nil {
		return nil, err
	}

	// アクティブユーザー数
	err = r.db.Model(&models.User{}).
		Where("is_active = ? AND deleted_at IS NULL", true).
		Count(&stats.ActiveUsers).Error
	if err != nil {
		return nil, err
	}

	// 非アクティブユーザー数
	err = r.db.Model(&models.User{}).
		Where("is_active = ? AND deleted_at IS NULL", false).
		Count(&stats.InactiveUsers).Error
	if err != nil {
		return nil, err
	}

	// 今日のログイン数（last_login_atが今日のユーザー）
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	err = r.db.Model(&models.User{}).
		Where("last_login_at >= ? AND last_login_at < ? AND deleted_at IS NULL", today, tomorrow).
		Count(&stats.TodayLogins).Error
	if err != nil {
		return nil, err
	}

	// アクティブセッション数
	err = r.db.Model(&models.UserSession{}).
		Where("is_revoked = ? AND expires_at > ? AND deleted_at IS NULL", false, time.Now()).
		Count(&stats.ActiveSessions).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// SearchUsers はユーザーを検索
func (r *UserRepository) SearchUsers(query string, offset, limit int) ([]models.UserListItem, error) {
	var users []models.UserListItem

	searchPattern := "%" + query + "%"
	err := r.db.Table("users").
		Select("id, username, email, first_name, last_name, is_active, last_login_at, created_at").
		Where("deleted_at IS NULL").
		Where("username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, err
}

// GetTotalUsersCount は総ユーザー数を取得
func (r *UserRepository) GetTotalUsersCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

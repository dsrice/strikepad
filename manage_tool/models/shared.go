package models

import (
	"time"

	"gorm.io/gorm"
)

// メインサービスのモデルを再利用するため、同じ構造体を定義
// 将来的にはメインサービスのモデルを直接importする方針

// User はユーザーモデル
type User struct {
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Password    string         `gorm:"not null" json:"-"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	GoogleID    string         `gorm:"index" json:"google_id,omitempty"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Picture     string         `json:"picture,omitempty"`
	Username    string         `gorm:"uniqueIndex;not null" json:"username"`
	ID          uint           `gorm:"primarykey" json:"id"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
}

// UserSession はユーザーセッションモデル
type UserSession struct {
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ExpiresAt    time.Time      `gorm:"not null" json:"expires_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	RefreshToken string         `gorm:"uniqueIndex;not null" json:"-"`
	IPAddress    string         `json:"ip_address,omitempty"`
	UserAgent    string         `json:"user_agent,omitempty"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	IsRevoked    bool           `gorm:"default:false" json:"is_revoked"`
}

// TableName はテーブル名を指定（メインサービスと同じテーブルを使用）
func (User) TableName() string {
	return "users"
}

func (UserSession) TableName() string {
	return "user_sessions"
}

// UserStats は管理ツール用のユーザー統計
type UserStats struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	InactiveUsers  int64 `json:"inactive_users"`
	TodayLogins    int64 `json:"today_logins"`
	ActiveSessions int64 `json:"active_sessions"`
}

// UserListItem はユーザー一覧表示用の軽量構造体
type UserListItem struct {
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	ID          uint       `json:"id"`
	IsActive    bool       `json:"is_active"`
}

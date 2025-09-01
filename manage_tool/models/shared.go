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

// AdminUser は管理者ユーザーモデル
type AdminUser struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	LoginID   string         `gorm:"size:30;not null" json:"login_id"`
	Password  string         `gorm:"size:200;not null" json:"-"`
	Name      string         `gorm:"size:30;not null" json:"name"`
	ID        uint           `gorm:"primarykey" json:"id"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
}

// TableName はテーブル名を指定
func (AdminUser) TableName() string {
	return "admin_users"
}

// Maker はメーカーモデル
type Maker struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	LogoFile  string         `gorm:"size:255" json:"logo_file,omitempty"`
	ID        uint           `gorm:"primarykey" json:"id"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
}

// TableName はテーブル名を指定
func (Maker) TableName() string {
	return "makers"
}

// Cover はカバーモデル
type Cover struct {
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Maker        Maker          `gorm:"foreignKey:MakerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name         string         `gorm:"size:200;not null" json:"name"`
	ID           uint           `gorm:"primarykey" json:"id"`
	MaterialType int            `gorm:"not null" json:"material_type"`
	MakerID      uint           `gorm:"not null;index" json:"maker_id"`
	Rank         int            `gorm:"not null" json:"rank"`
	IsDeleted    bool           `gorm:"default:false" json:"is_deleted"`
}

// TableName はテーブル名を指定
func (Cover) TableName() string {
	return "covers"
}

// Core はコアモデル
type Core struct {
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	InitDiff     *float32       `json:"init_diff,omitempty"`
	Maker        Maker          `gorm:"foreignKey:MakerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name         string         `gorm:"size:200;not null" json:"name"`
	ID           uint           `gorm:"primarykey" json:"id"`
	MakerID      uint           `gorm:"not null;index" json:"maker_id"`
	RG           float32        `gorm:"not null" json:"rg"`
	DeltaRG      float32        `gorm:"not null" json:"delta_rg"`
	SymmetryFlag bool           `gorm:"not null" json:"symmetry_flag"`
	IsDeleted    bool           `gorm:"default:false" json:"is_deleted"`
}

// TableName はテーブル名を指定
func (Core) TableName() string {
	return "cores"
}

// Ball はボールモデル
type Ball struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Maker     Maker          `gorm:"foreignKey:MakerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Core      Core           `gorm:"foreignKey:CoreID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Cover     Cover          `gorm:"foreignKey:CoverID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	URL       string         `gorm:"size:200;not null" json:"url"`
	ID        uint           `gorm:"primarykey" json:"id"`
	MakerID   uint           `gorm:"not null" json:"maker_id"`
	CoreID    uint           `gorm:"not null" json:"core_id"`
	CoverID   uint           `gorm:"not null" json:"cover_id"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
}

// TableName はテーブル名を指定
func (Ball) TableName() string {
	return "balls"
}

// MakerListItem はメーカー一覧表示用の軽量構造体
type MakerListItem struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name       string    `json:"name"`
	ID         uint      `json:"id"`
	BallCount  int64     `json:"ball_count"`
	CoverCount int64     `json:"cover_count"`
	CoreCount  int64     `json:"core_count"`
}

// MakerStats はメーカー統計用構造体
type MakerStats struct {
	TotalMakers  int64 `json:"total_makers"`
	ActiveMakers int64 `json:"active_makers"`
	TotalBalls   int64 `json:"total_balls"`
	TotalCovers  int64 `json:"total_covers"`
	TotalCores   int64 `json:"total_cores"`
}
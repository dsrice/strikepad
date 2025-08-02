package model

import (
	"time"
)

type User struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	EmailVerified  bool       `gorm:"column:email_verified;default:false;not null" json:"email_verified"`
	IsDeleted      bool       `gorm:"column:is_deleted;default:false;not null" json:"-"`
	CreatedAt      time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"-"`
	ProviderUserID *string    `gorm:"column:provider_user_id;size:255" json:"provider_user_id,omitempty"`
	Email          *string    `gorm:"column:email;size:255" json:"email,omitempty"`
	PasswordHash   *string    `gorm:"column:password_hash;size:255" json:"-"`
	ProviderType   string     `gorm:"column:provider_type;size:20;not null" json:"provider_type"`
	DisplayName    string     `gorm:"column:display_name;size:100;not null" json:"display_name"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

package model

import (
	"time"
)

type User struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	ProviderType   string     `gorm:"column:provider_type;size:20;not null" json:"provider_type"`
	ProviderUserID *string    `gorm:"column:provider_user_id;size:255" json:"provider_user_id,omitempty"`
	Email          *string    `gorm:"column:email;size:255" json:"email,omitempty"`
	DisplayName    string     `gorm:"column:display_name;size:100;not null" json:"display_name"`
	PasswordHash   *string    `gorm:"column:password_hash;size:255" json:"-"`
	EmailVerified  bool       `gorm:"column:email_verified;default:false;not null" json:"email_verified"`
	CreatedAt      time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null" json:"updated_at"`
	IsDeleted      bool       `gorm:"column:is_deleted;default:false;not null" json:"-"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

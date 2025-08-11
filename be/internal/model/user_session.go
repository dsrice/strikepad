package model

import (
	"time"

	"gorm.io/gorm"
)

// UserSession represents a user session with tokens
type UserSession struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	UserID                uint           `gorm:"not null;index" json:"user_id"`
	AccessToken           string         `gorm:"type:text;not null" json:"access_token"`
	RefreshToken          string         `gorm:"column:refresh_token;type:text" json:"refresh_token"`
	AccessTokenExpiresAt  time.Time      `gorm:"not null" json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time      `gorm:"not null" json:"refresh_token_expires_at"`
	CreatedAt             time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsDeleted             bool           `gorm:"default:false" json:"is_deleted"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at,omitempty"`

	// Relationship
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

// TableName returns the table name for GORM
func (UserSession) TableName() string {
	return "user_sessions"
}

// IsAccessTokenValid checks if the access token is still valid
func (us *UserSession) IsAccessTokenValid() bool {
	return time.Now().Before(us.AccessTokenExpiresAt) && !us.IsDeleted
}

// IsRefreshTokenValid checks if the refresh token is still valid
func (us *UserSession) IsRefreshTokenValid() bool {
	return time.Now().Before(us.RefreshTokenExpiresAt) && !us.IsDeleted
}

// Invalidate marks the session as deleted
func (us *UserSession) Invalidate() {
	us.IsDeleted = true
	us.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
}
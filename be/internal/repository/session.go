package repository

import (
	"fmt"
	"time"

	"strikepad-backend/internal/model"

	"gorm.io/gorm"
)

// SessionRepository handles database operations for user sessions
type SessionRepository struct {
	db *gorm.DB
}

// SessionRepositoryInterface defines the interface for session repository
type SessionRepositoryInterface interface {
	Create(session *model.UserSession) error
	FindByAccessToken(accessToken string) (*model.UserSession, error)
	FindByRefreshToken(refreshToken string) (*model.UserSession, error)
	FindActiveByUserID(userID uint) ([]*model.UserSession, error)
	Update(session *model.UserSession) error
	InvalidateByUserID(userID uint) error
	InvalidateExpiredSessions() error
	Delete(sessionID uint) error
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB) SessionRepositoryInterface {
	return &SessionRepository{
		db: db,
	}
}

// Create creates a new user session
func (r *SessionRepository) Create(session *model.UserSession) error {
	if err := r.db.Create(session).Error; err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

// FindByAccessToken finds a session by access token
func (r *SessionRepository) FindByAccessToken(accessToken string) (*model.UserSession, error) {
	var session model.UserSession
	err := r.db.Where("access_token = ? AND is_deleted = false", accessToken).
		Preload("User").
		First(&session).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to find session by access token: %w", err)
	}

	return &session, nil
}

// FindByRefreshToken finds a session by refresh token
func (r *SessionRepository) FindByRefreshToken(refreshToken string) (*model.UserSession, error) {
	var session model.UserSession
	err := r.db.Where("refresh_token = ? AND is_deleted = false", refreshToken).
		Preload("User").
		First(&session).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to find session by refresh token: %w", err)
	}

	return &session, nil
}

// FindActiveByUserID finds all active sessions for a user
func (r *SessionRepository) FindActiveByUserID(userID uint) ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	err := r.db.Where("user_id = ? AND is_deleted = false AND access_token_expires_at > ?",
		userID, time.Now()).
		Order("created_at DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find active sessions: %w", err)
	}

	return sessions, nil
}

// Update updates a session
func (r *SessionRepository) Update(session *model.UserSession) error {
	if err := r.db.Save(session).Error; err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

// InvalidateByUserID invalidates all sessions for a specific user
func (r *SessionRepository) InvalidateByUserID(userID uint) error {
	err := r.db.Model(&model.UserSession{}).
		Where("user_id = ? AND is_deleted = false", userID).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to invalidate sessions for user %d: %w", userID, err)
	}

	return nil
}

// InvalidateExpiredSessions marks expired sessions as deleted
func (r *SessionRepository) InvalidateExpiredSessions() error {
	now := time.Now()

	// Invalidate sessions where both tokens are expired
	err := r.db.Model(&model.UserSession{}).
		Where("is_deleted = false AND refresh_token_expires_at < ?", now).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": now,
			"updated_at": now,
		}).Error

	if err != nil {
		return fmt.Errorf("failed to invalidate expired sessions: %w", err)
	}

	return nil
}

// Delete permanently deletes a session
func (r *SessionRepository) Delete(sessionID uint) error {
	if err := r.db.Delete(&model.UserSession{}, sessionID).Error; err != nil {
		return fmt.Errorf("failed to delete session %d: %w", sessionID, err)
	}
	return nil
}

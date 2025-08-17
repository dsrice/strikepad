package service

import (
	"fmt"
	"log/slog"
	"time"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/model"
	"strikepad-backend/internal/repository"
)

// SessionService handles session-related business logic
type SessionService struct {
	sessionRepo repository.SessionRepositoryInterface
	jwtService  *auth.JWTService
}

// SessionServiceInterface defines the interface for session service
type SessionServiceInterface interface {
	CreateSession(userID uint) (*auth.TokenPair, error)
	ValidateAccessToken(token string) (*model.UserSession, error)
	RefreshToken(refreshToken string) (*auth.TokenPair, error)
	RefreshSession(accessToken, refreshToken string) (*auth.TokenPair, error)
	InvalidateSession(accessToken string) error
	InvalidateAllUserSessions(userID uint) error
	Logout(userID uint, accessToken string) error
	CleanupExpiredSessions() error
}

// NewSessionService creates a new session service
func NewSessionService(
	sessionRepo repository.SessionRepositoryInterface,
	jwtService *auth.JWTService,
) SessionServiceInterface {
	return &SessionService{
		sessionRepo: sessionRepo,
		jwtService:  jwtService,
	}
}

// CreateSession creates a new session with token pair
func (s *SessionService) CreateSession(userID uint) (*auth.TokenPair, error) {
	// Generate token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// Create session record
	session := &model.UserSession{
		UserID:                userID,
		AccessToken:           tokenPair.AccessToken,
		RefreshToken:          tokenPair.RefreshToken,
		AccessTokenExpiresAt:  tokenPair.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: tokenPair.RefreshTokenExpiresAt,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		IsDeleted:             false,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	slog.Info("Session created successfully", "user_id", userID, "session_id", session.ID)
	return tokenPair, nil
}

// ValidateAccessToken validates an access token and returns the session
func (s *SessionService) ValidateAccessToken(token string) (*model.UserSession, error) {
	// Validate JWT token
	claims, err := s.jwtService.ValidateAccessToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	// Find session in database
	session, err := s.sessionRepo.FindByAccessToken(token)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Check if session is still valid
	if !session.IsAccessTokenValid() {
		return nil, fmt.Errorf("session is expired or invalidated")
	}

	// Verify user ID matches
	if session.UserID != claims.UserID {
		return nil, fmt.Errorf("token user ID mismatch")
	}

	return session, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *SessionService) RefreshToken(refreshToken string) (*auth.TokenPair, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Find session in database
	session, err := s.sessionRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Check if refresh token is still valid
	if !session.IsRefreshTokenValid() {
		return nil, fmt.Errorf("refresh token is expired or invalidated")
	}

	// Verify user ID matches
	if session.UserID != claims.UserID {
		return nil, fmt.Errorf("token user ID mismatch")
	}

	// Generate new token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token pair: %w", err)
	}

	// Update session with new tokens
	session.AccessToken = tokenPair.AccessToken
	session.RefreshToken = tokenPair.RefreshToken
	session.AccessTokenExpiresAt = tokenPair.AccessTokenExpiresAt
	session.RefreshTokenExpiresAt = tokenPair.RefreshTokenExpiresAt
	session.UpdatedAt = time.Now()

	if err := s.sessionRepo.Update(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	slog.Info("Token refreshed successfully", "user_id", claims.UserID, "session_id", session.ID)
	return tokenPair, nil
}

// RefreshSession refreshes tokens using both access and refresh tokens
func (s *SessionService) RefreshSession(accessToken, refreshToken string) (*auth.TokenPair, error) {
	// Validate refresh token first
	refreshClaims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Find session by refresh token
	session, err := s.sessionRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Check if refresh token is still valid
	if !session.IsRefreshTokenValid() {
		return nil, fmt.Errorf("refresh token is expired or invalidated")
	}

	// Verify the provided access token matches the session
	if session.AccessToken != accessToken {
		return nil, fmt.Errorf("access token does not match session")
	}

	// Verify user ID matches
	if session.UserID != refreshClaims.UserID {
		return nil, fmt.Errorf("token user ID mismatch")
	}

	// Generate new token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(refreshClaims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token pair: %w", err)
	}

	// Update session with new tokens
	session.AccessToken = tokenPair.AccessToken
	session.RefreshToken = tokenPair.RefreshToken
	session.AccessTokenExpiresAt = tokenPair.AccessTokenExpiresAt
	session.RefreshTokenExpiresAt = tokenPair.RefreshTokenExpiresAt
	session.UpdatedAt = time.Now()

	if err := s.sessionRepo.Update(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	slog.Info("Session refreshed successfully", "user_id", refreshClaims.UserID, "session_id", session.ID)
	return tokenPair, nil
}

// InvalidateSession invalidates a session by access token
func (s *SessionService) InvalidateSession(accessToken string) error {
	session, err := s.sessionRepo.FindByAccessToken(accessToken)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	session.Invalidate()
	if err := s.sessionRepo.Update(session); err != nil {
		return fmt.Errorf("failed to invalidate session: %w", err)
	}

	slog.Info("Session invalidated successfully", "user_id", session.UserID, "session_id", session.ID)
	return nil
}

// InvalidateAllUserSessions invalidates all sessions for a specific user
func (s *SessionService) InvalidateAllUserSessions(userID uint) error {
	if err := s.sessionRepo.InvalidateByUserID(userID); err != nil {
		return fmt.Errorf("failed to invalidate all user sessions: %w", err)
	}

	slog.Info("All user sessions invalidated", "user_id", userID)
	return nil
}

// Logout handles user logout by invalidating the specific session
func (s *SessionService) Logout(userID uint, accessToken string) error {
	// Find session by access token
	session, err := s.sessionRepo.FindByAccessToken(accessToken)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Verify the session belongs to the user
	if session.UserID != userID {
		return fmt.Errorf("session does not belong to user")
	}

	// Invalidate the session
	session.Invalidate()
	if err := s.sessionRepo.Update(session); err != nil {
		return fmt.Errorf("failed to logout session: %w", err)
	}

	slog.Info("User logged out successfully", "user_id", userID, "session_id", session.ID)
	return nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *SessionService) CleanupExpiredSessions() error {
	if err := s.sessionRepo.InvalidateExpiredSessions(); err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	slog.Info("Expired sessions cleaned up successfully")
	return nil
}

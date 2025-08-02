package auth

import (
	"regexp"
	"strings"
)

// Email validation regex pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	// Trim whitespace
	email = strings.TrimSpace(email)

	// Check email format
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

// NormalizeEmail normalizes email to lowercase and trims whitespace
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

package auth

import "errors"

var (
	// Password validation errors
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong  = errors.New("password must be at most 128 characters long")

	// Email validation errors
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrEmailRequired = errors.New("email is required")

	// User creation errors
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrUserNotFound      = errors.New("user not found")

	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid email or password")
)

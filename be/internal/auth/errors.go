package auth

import "errors"

var (
	// ErrPasswordTooShort is returned when password is shorter than minimum length
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	// ErrPasswordTooLong is returned when password exceeds maximum length
	ErrPasswordTooLong = errors.New("password must be at most 128 characters long")

	// ErrInvalidEmail is returned when email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")
	// ErrEmailRequired is returned when email is missing
	ErrEmailRequired = errors.New("email is required")

	// ErrUserAlreadyExists is returned when attempting to create a user that already exists
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	// ErrUserNotFound is returned when requested user does not exist
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials is returned when login credentials are incorrect
	ErrInvalidCredentials = errors.New("invalid email or password")
)

package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PasswordTestSuite struct {
	suite.Suite
}

func (suite *PasswordTestSuite) TestHashPassword() {
	password := "testPassword123"

	hash, err := HashPassword(password)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hash)
	assert.NotEqual(suite.T(), password, hash)

	// Hash should start with bcrypt prefix
	assert.True(suite.T(), strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$"))
}

func (suite *PasswordTestSuite) TestHashPasswordEmptyString() {
	hash, err := HashPassword("")
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hash)
}

func (suite *PasswordTestSuite) TestCheckPasswordHashValid() {
	password := "testPassword123"

	hash, err := HashPassword(password)
	assert.NoError(suite.T(), err)

	// Check with correct password
	isValid := CheckPasswordHash(password, hash)
	assert.True(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestCheckPasswordHashInvalid() {
	password := "testPassword123"
	wrongPassword := "wrongPassword456"

	hash, err := HashPassword(password)
	assert.NoError(suite.T(), err)

	// Check with wrong password
	isValid := CheckPasswordHash(wrongPassword, hash)
	assert.False(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestCheckPasswordHashInvalidHash() {
	password := "testPassword123"
	invalidHash := "invalid_hash"

	// Check with invalid hash format
	isValid := CheckPasswordHash(password, invalidHash)
	assert.False(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestCheckPasswordHashEmptyPassword() {
	hash, err := HashPassword("testPassword123")
	assert.NoError(suite.T(), err)

	// Check with empty password
	isValid := CheckPasswordHash("", hash)
	assert.False(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestCheckPasswordHashEmptyHash() {
	// Check with empty hash
	isValid := CheckPasswordHash("testPassword123", "")
	assert.False(suite.T(), isValid)
}

func (suite *PasswordTestSuite) TestValidatePasswordValid() {
	testCases := []struct {
		name     string
		password string
	}{
		{"Minimum length", "Password1!"},
		{"Medium length", "MySecurePassword123!"},
		{"Maximum length", "Password123!" + strings.Repeat("a", 116)}, // 128 chars total
		{"Special characters", "Test@#$%^&*()123A"},
		{"Unicode characters", "Пароль123!"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)
			assert.NoError(t, err, "Password should be valid: %s", tc.password)
		})
	}
}

func (suite *PasswordTestSuite) TestValidatePasswordTooShort() {
	shortPasswords := []string{
		"",
		"a",
		"Pass1!",  // 6 chars
		"Test12!", // 7 chars
	}

	for _, password := range shortPasswords {
		err := ValidatePassword(password)
		assert.Equal(suite.T(), ErrPasswordTooShort, err, "Password should be too short: %s", password)
	}
}

func (suite *PasswordTestSuite) TestValidatePasswordTooLong() {
	// Create password longer than 128 characters
	longPassword := "Password123!" + strings.Repeat("a", 120) // 132 chars total

	err := ValidatePassword(longPassword)
	assert.Equal(suite.T(), ErrPasswordTooLong, err)
}

func (suite *PasswordTestSuite) TestValidatePasswordExactLimits() {
	// Test exactly 8 characters (minimum)
	minPassword := "Pass123!"
	assert.Len(suite.T(), minPassword, 8)
	err := ValidatePassword(minPassword)
	assert.NoError(suite.T(), err)

	// Test exactly 128 characters (maximum)
	maxPassword := "Pass123!" + strings.Repeat("a", 120) // 128 chars total
	assert.Len(suite.T(), maxPassword, 128)
	err = ValidatePassword(maxPassword)
	assert.NoError(suite.T(), err)
}

func (suite *PasswordTestSuite) TestHashPasswordConsistency() {
	password := "testPassword123"

	// Hash the same password multiple times
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)

	// Hashes should be different (due to salt)
	assert.NotEqual(suite.T(), hash1, hash2)

	// But both should verify correctly
	assert.True(suite.T(), CheckPasswordHash(password, hash1))
	assert.True(suite.T(), CheckPasswordHash(password, hash2))
}

func (suite *PasswordTestSuite) TestPasswordWorkflow() {
	// Test a complete password workflow
	originalPassword := "MySecurePassword123!"

	// 1. Validate password
	err := ValidatePassword(originalPassword)
	assert.NoError(suite.T(), err)

	// 2. Hash password
	hash, err := HashPassword(originalPassword)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hash)

	// 3. Verify correct password
	isValid := CheckPasswordHash(originalPassword, hash)
	assert.True(suite.T(), isValid)

	// 4. Verify wrong password fails
	wrongPassword := "WrongPassword456!"
	isValid = CheckPasswordHash(wrongPassword, hash)
	assert.False(suite.T(), isValid)
}

func TestPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}

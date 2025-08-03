package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const testPasswordConst = "testPasswordConst123"

type PasswordTestSuite struct {
	suite.Suite
}

func (suite *PasswordTestSuite) TestHashPassword() {
	testCases := []struct {
		name           string
		password       string
		expectError    bool
		validateFormat bool
		validateUnique bool
	}{
		{
			name:           "valid password",
			password:       testPasswordConst,
			expectError:    false,
			validateFormat: true,
			validateUnique: true,
		},
		{
			name:           "empty string",
			password:       "",
			expectError:    false,
			validateFormat: true,
			validateUnique: false,
		},
		{
			name:           "long password",
			password:       "very_long_password_" + strings.Repeat("a", 30), // 49 chars total, under 72 byte limit
			expectError:    false,
			validateFormat: true,
			validateUnique: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			hash, err := HashPassword(tc.password)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tc.password, hash)

			if tc.validateFormat {
				// Hash should start with bcrypt prefix
				assert.True(t, strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$"))
			}

			if tc.validateUnique {
				// Generate another hash and ensure they're different
				hash2, err2 := HashPassword(tc.password)
				assert.NoError(t, err2)
				assert.NotEqual(t, hash, hash2, "Two hashes of the same password should be different due to salt")
			}
		})
	}
}

func (suite *PasswordTestSuite) TestCheckPasswordHash() {
	// Generate a valid hash for testing
	validPassword := testPasswordConst
	validHash, err := HashPassword(validPassword)
	assert.NoError(suite.T(), err)

	testCases := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "correct password",
			password: validPassword,
			hash:     validHash,
			expected: true,
		},
		{
			name:     "wrong password",
			password: "wrongPassword456",
			hash:     validHash,
			expected: false,
		},
		{
			name:     "invalid hash format",
			password: validPassword,
			hash:     "invalid_hash",
			expected: false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     validHash,
			expected: false,
		},
		{
			name:     "empty hash",
			password: validPassword,
			hash:     "",
			expected: false,
		},
		{
			name:     "both empty",
			password: "",
			hash:     "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			isValid := CheckPasswordHash(tc.password, tc.hash)
			assert.Equal(t, tc.expected, isValid)
		})
	}
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

func (suite *PasswordTestSuite) TestValidatePasswordLength() {
	testCases := []struct {
		name        string
		password    string
		expectedErr error
		description string
	}{
		// Too short passwords
		{"empty password", "", ErrPasswordTooShort, "0 chars"},
		{"single character", "a", ErrPasswordTooShort, "1 char"},
		{"six characters", "Pass1!", ErrPasswordTooShort, "6 chars"},
		{"seven characters", "Test12!", ErrPasswordTooShort, "7 chars"},

		// Valid length passwords
		{"minimum valid", "Pass123!", nil, "8 chars (minimum)"},
		{"maximum valid", "Pass123!" + strings.Repeat("a", 120), nil, "128 chars (maximum)"},
		{"medium length", "MySecurePassword123!", nil, "20 chars"},

		// Too long passwords
		{"too long", "Password123!" + strings.Repeat("a", 120), ErrPasswordTooLong, "132 chars"},
		{"way too long", strings.Repeat("a", 200), ErrPasswordTooLong, "200 chars"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err, "Password validation failed for %s", tc.description)
			} else {
				assert.NoError(t, err, "Password should be valid for %s", tc.description)
			}
		})
	}
}

func (suite *PasswordTestSuite) TestHashPasswordConsistency() {
	testCases := []struct {
		name     string
		password string
	}{
		{"standard password", testPasswordConst},
		{"empty password", ""},
		{"special characters", "P@ssw0rd!#$%"},
		{"unicode password", "パスワード123!"},
		{"long password", "MyVeryLongPassword" + strings.Repeat("a", 45)}, // 65 chars total, under 72 byte limit
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Hash the same password multiple times
			hash1, err1 := HashPassword(tc.password)
			hash2, err2 := HashPassword(tc.password)

			assert.NoError(t, err1)
			assert.NoError(t, err2)

			// Hashes should be different (due to salt)
			assert.NotEqual(t, hash1, hash2, "Two hashes of the same password should be different due to salt")

			// But both should verify correctly
			assert.True(t, CheckPasswordHash(tc.password, hash1), "First hash should verify correctly")
			assert.True(t, CheckPasswordHash(tc.password, hash2), "Second hash should verify correctly")

			// Cross-verification should also work
			assert.True(t, CheckPasswordHash(tc.password, hash1), "Password should verify against first hash")
			assert.True(t, CheckPasswordHash(tc.password, hash2), "Password should verify against second hash")
		})
	}
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
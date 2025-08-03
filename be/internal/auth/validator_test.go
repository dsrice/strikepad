package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthValidatorTestSuite struct {
	suite.Suite
}

func (suite *AuthValidatorTestSuite) TestValidateEmail() {
	testCases := []struct {
		name      string
		email     string
		expectErr error
	}{
		// Valid emails
		{"valid basic email", "test@example.com", nil},
		{"valid with dots", "user.name@domain.co.uk", nil},
		{"valid with plus", "user+tag@example.org", nil},
		{"valid with underscore", "user_name@example-domain.com", nil},
		{"valid with numbers", "123@numbers.com", nil},
		{"valid subdomain", "test@sub.domain.example.com", nil},
		{"valid long email", "very.long.email.address@very.long.domain.name.example.com", nil},
		
		// Invalid emails
		{"empty email", "", ErrEmailRequired},
		{"no @ symbol", "invalid", ErrInvalidEmail},
		{"no local part", "@example.com", ErrInvalidEmail},
		{"no domain", "test@", ErrInvalidEmail},
		{"no TLD", "test@example", ErrInvalidEmail},
		{"space in local", "te st@example.com", ErrInvalidEmail},
		{"space in domain", "test@exam ple.com", ErrInvalidEmail},
		{"tab in domain", "test@exam\tple.com", ErrInvalidEmail},
		{"newline in domain", "test@exam\nple.com", ErrInvalidEmail},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.email)
			if tc.expectErr != nil {
				assert.Equal(t, tc.expectErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *AuthValidatorTestSuite) TestValidateEmailSpecialCases() {
	testCases := []struct {
		name  string
		email string
		valid bool
	}{
		{"long valid email", "test.user.with.long.name@example.com", true},
		{"plus tag", "test+tag@example.com", true},
		{"dash tag", "test-tag@example.com", true},
		{"underscore tag", "test_tag@example.com", true},
		{"dot tag", "test.tag@example.com", true},
		{"numbers only", "123456@example.com", true},
		{"space before", "test @example.com", false},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.email)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func (suite *AuthValidatorTestSuite) TestNormalizeEmail() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic normalization
		{"uppercase to lowercase", "Test@Example.COM", "test@example.com"},
		{"mixed case", "USER@DOMAIN.ORG", "user@domain.org"},
		{"already lowercase", "already@lowercase.com", "already@lowercase.com"},
		{"numbers", "123@Numbers.Com", "123@numbers.com"},
		
		// Trimming whitespace
		{"leading and trailing spaces", "  test@example.com  ", "test@example.com"},
		{"tabs", "\tuser@domain.org\t", "user@domain.org"},
		{"newlines", "\nuser@domain.org\n", "user@domain.org"},
		{"mixed whitespace", " Test@Example.COM ", "test@example.com"},
		{"multiple spaces", "   UPPER@CASE.COM   ", "upper@case.com"},
		
		// Preserving structure
		{"plus tag", "User+Tag@Example.Com", "user+tag@example.com"},
		{"dots", "User.Name@Sub.Domain.Com", "user.name@sub.domain.com"},
		{"underscores", "User_Name@Example-Domain.Org", "user_name@example-domain.org"},
		{"numbers with dots", "123.456@Numbers.Co.UK", "123.456@numbers.co.uk"},
		
		// Edge cases
		{"empty string", "", ""},
		{"whitespace only", "   ", ""},
		{"tabs only", "\t\t", ""},
		{"newlines only", "\n\n", ""},
		{"mixed whitespace only", " \t\n ", ""},
		
		// Unicode
		{"unicode characters", "TëSt@Example.com", "tëst@example.com"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			result := NormalizeEmail(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func (suite *AuthValidatorTestSuite) TestEmailValidationWorkflow() {
	testCases := []struct {
		name           string
		inputEmail     string
		expectedNorm   string
		shouldValidate bool
	}{
		{
			name:           "complete workflow",
			inputEmail:     "  User.Name+Tag@Example.COM  ",
			expectedNorm:   "user.name+tag@example.com",
			shouldValidate: true,
		},
		{
			name:           "already normalized",
			inputEmail:     "test@example.com",
			expectedNorm:   "test@example.com",
			shouldValidate: true,
		},
		{
			name:           "invalid email",
			inputEmail:     "invalid-email",
			expectedNorm:   "invalid-email",
			shouldValidate: false,
		},
		{
			name:           "empty email",
			inputEmail:     "   ",
			expectedNorm:   "",
			shouldValidate: false,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Test normalization
			normalized := NormalizeEmail(tc.inputEmail)
			assert.Equal(t, tc.expectedNorm, normalized)

			// Test validation of original
			err := ValidateEmail(tc.inputEmail)
			if tc.shouldValidate {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			// Test validation of normalized (if not empty)
			if normalized != "" {
				err = ValidateEmail(normalized)
				if tc.shouldValidate {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			}
		})
	}
}


func TestAuthValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(AuthValidatorTestSuite))
}

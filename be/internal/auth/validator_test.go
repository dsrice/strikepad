package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthValidatorTestSuite struct {
	suite.Suite
}

func (suite *AuthValidatorTestSuite) TestValidateEmailValid() {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"user+tag@example.org",
		"user_name@example-domain.com",
		"123@numbers.com",
		"test@sub.domain.example.com",
		"very.long.email.address@very.long.domain.name.example.com",
	}

	for _, email := range validEmails {
		suite.T().Run(email, func(t *testing.T) {
			err := ValidateEmail(email)
			assert.NoError(t, err, "Email should be valid: %s", email)
		})
	}
}

func (suite *AuthValidatorTestSuite) TestValidateEmailInvalid() {
	invalidEmails := []string{
		"",                   // Empty - returns ErrEmailRequired
		"invalid",            // No @ symbol
		"@example.com",       // No local part
		"test@",              // No domain
		"test@example",       // No TLD
		"test @example.com",  // Space in local part (before trim)
		"test@exam ple.com",  // Space in domain
		"te st@example.com",  // Space in local part
		"test@exam\tple.com", // Tab in domain
		"test@exam\nple.com", // Newline in domain
	}

	for _, email := range invalidEmails {
		suite.T().Run("invalid_"+email, func(t *testing.T) {
			err := ValidateEmail(email)
			if email == "" {
				assert.Equal(t, ErrEmailRequired, err, "Empty email should return ErrEmailRequired")
			} else {
				assert.Equal(t, ErrInvalidEmail, err, "Email should be invalid: %s", email)
			}
		})
	}
}

func (suite *AuthValidatorTestSuite) TestValidateEmailEdgeCases() {
	// Test a reasonably long email that should be valid
	localPart := "test.user.with.long.name"
	domain := "example.com"
	longEmail := localPart + "@" + domain

	// This should be valid
	err := ValidateEmail(longEmail)
	assert.NoError(suite.T(), err)

	// Test with special characters that should be valid
	specialEmails := []string{
		"test+tag@example.com",
		"test-tag@example.com",
		"test_tag@example.com",
		"test.tag@example.com",
		"123456@example.com",
	}

	for _, email := range specialEmails {
		err := ValidateEmail(email)
		assert.NoError(suite.T(), err, "Special email should be valid: %s", email)
	}
}

func (suite *AuthValidatorTestSuite) TestNormalizeEmailBasic() {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Test@Example.COM", "test@example.com"},
		{"USER@DOMAIN.ORG", "user@domain.org"},
		{"Mixed.Case@Email.Co.UK", "mixed.case@email.co.uk"},
		{"already@lowercase.com", "already@lowercase.com"},
		{"123@Numbers.Com", "123@numbers.com"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.input, func(t *testing.T) {
			result := NormalizeEmail(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func (suite *AuthValidatorTestSuite) TestNormalizeEmailTrimsSpaces() {
	testCases := []struct {
		input    string
		expected string
	}{
		{"  test@example.com  ", "test@example.com"},
		{"\tuser@domain.org\t", "user@domain.org"},
		{"\nuser@domain.org\n", "user@domain.org"},
		{" Test@Example.COM ", "test@example.com"},
		{"   UPPER@CASE.COM   ", "upper@case.com"},
	}

	for _, tc := range testCases {
		suite.T().Run("trim_"+tc.input, func(t *testing.T) {
			result := NormalizeEmail(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func (suite *AuthValidatorTestSuite) TestNormalizeEmailPreservesStructure() {
	testCases := []string{
		"user+tag@example.com",
		"user.name@sub.domain.com",
		"user_name@example-domain.org",
		"123.456@numbers.co.uk",
	}

	for _, email := range testCases {
		result := NormalizeEmail(email)

		// Should be lowercase
		assert.Equal(suite.T(), email, result)

		// Should preserve the @ symbol position
		inputAtPos := -1
		resultAtPos := -1
		for i, char := range email {
			if char == '@' {
				inputAtPos = i
				break
			}
		}
		for i, char := range result {
			if char == '@' {
				resultAtPos = i
				break
			}
		}
		assert.Equal(suite.T(), inputAtPos, resultAtPos)
	}
}

func (suite *AuthValidatorTestSuite) TestNormalizeEmailEmpty() {
	result := NormalizeEmail("")
	assert.Equal(suite.T(), "", result)
}

func (suite *AuthValidatorTestSuite) TestNormalizeEmailWhitespaceOnly() {
	whitespaceInputs := []string{
		"   ",
		"\t\t",
		"\n\n",
		" \t\n ",
	}

	for _, input := range whitespaceInputs {
		result := NormalizeEmail(input)
		assert.Equal(suite.T(), "", result, "Whitespace-only input should result in empty string")
	}
}

func (suite *AuthValidatorTestSuite) TestEmailValidationAndNormalizationWorkflow() {
	// Test a complete email validation and normalization workflow
	inputEmail := "  User.Name+Tag@Example.COM  "

	// 1. Normalize email
	normalized := NormalizeEmail(inputEmail)
	expected := "user.name+tag@example.com"
	assert.Equal(suite.T(), expected, normalized)

	// 2. Validate normalized email
	err := ValidateEmail(normalized)
	assert.NoError(suite.T(), err)

	// 3. Original email (with validation) should pass because ValidateEmail trims whitespace
	err = ValidateEmail(inputEmail)
	assert.NoError(suite.T(), err)

	// 4. But normalized email should be valid
	err = ValidateEmail(normalized)
	assert.NoError(suite.T(), err)
}

func (suite *AuthValidatorTestSuite) TestNormalizeEmailUnicode() {
	// Test with some basic unicode characters
	testCases := []struct {
		input    string
		expected string
	}{
		{"TëSt@Example.com", "tëst@example.com"},
		{"User@EXAMPLE.COM", "user@example.com"},
		{"Test@Example.COM", "test@example.com"},
	}

	for _, tc := range testCases {
		result := NormalizeEmail(tc.input)
		// Should at least not panic and return something
		assert.NotEmpty(suite.T(), result)
		// Should be lowercase of the input
		assert.Equal(suite.T(), tc.expected, result)
	}
}

func TestAuthValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(AuthValidatorTestSuite))
}

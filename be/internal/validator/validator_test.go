package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const testPasswordFieldConstConst = "password"

type ValidatorTestSuite struct {
	suite.Suite
	validator *Validator
}

func (suite *ValidatorTestSuite) SetupTest() {
	suite.validator = New()
}

// Test structs for validation
type TestUser struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:testPasswordFieldConst validate:"required,min=8,max=128,password_complex"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
	Age         int    `json:"age" validate:"gte=0,lte=150"`
}

type TestProduct struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Price       int    `json:"price" validate:"required,gt=0"`
	Category    string `json:"category" validate:"required,oneof=electronics clothing books"`
	Description string `json:"description" validate:"max=500"`
	Code        string `json:"code" validate:"required,alphanum,len=8"`
}

func (suite *ValidatorTestSuite) TestValidatePasswordComplexity() {
	testCases := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Valid complex password",
			password: "Password123!",
			expected: true,
		},
		{
			name:     "Valid complex password with symbols",
			password: "MyP@ssw0rd!",
			expected: true,
		},
		{
			name:     "Missing uppercase",
			password: "password123!",
			expected: false,
		},
		{
			name:     "Missing lowercase",
			password: "PASSWORD123!",
			expected: false,
		},
		{
			name:     "Missing symbol",
			password: "Password123",
			expected: false,
		},
		{
			name:     "Only lowercase",
			password: testPasswordFieldConst,
			expected: false,
		},
		{
			name:     "Only uppercase",
			password: "PASSWORD",
			expected: false,
		},
		{
			name:     "Only numbers",
			password: "12345678",
			expected: false,
		},
		{
			name:     "Empty password",
			password: "",
			expected: false, // Will fail on required, not password_complex
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			user := TestUser{
				Email:       "test@example.com",
				Password:    tc.password,
				DisplayName: "Test User",
				Age:         25,
			}

			err := suite.validator.Validate(&user)

			if tc.expected {
				// Should not have password_complex validation error
				if err != nil {
					ve, ok := err.(ValidationErrors)
					assert.True(t, ok)
					for _, validationErr := range ve.Errors {
						assert.NotEqual(t, "password_complex", validationErr.Tag, "Should not have password complexity error for valid password")
					}
				}
			} else {
				// Should have validation error (either password_complex or required)
				assert.Error(t, err)
				ve, ok := err.(ValidationErrors)
				assert.True(t, ok)

				// For empty password, it will fail on "required" first, not "password_complex"
				if tc.password == "" {
					hasRequiredError := false
					for _, validationErr := range ve.Errors {
						if validationErr.Tag == "required" && validationErr.Field == testPasswordFieldConst {
							hasRequiredError = true
							break
						}
					}
					assert.True(t, hasRequiredError, "Empty password should have required error")
				} else {
					// Non-empty but invalid passwords should have password_complex error
					hasPasswordComplexError := false
					for _, validationErr := range ve.Errors {
						if validationErr.Tag == "password_complex" {
							hasPasswordComplexError = true
							assert.Equal(t, testPasswordFieldConst, validationErr.Field)
							assert.Contains(t, validationErr.Message, "must contain at least one lowercase letter, one uppercase letter, and one symbol")
							break
						}
					}
					assert.True(t, hasPasswordComplexError, "Should have password complexity error")
				}
			}
		})
	}
}

func (suite *ValidatorTestSuite) TestValidateSuccess() {
	user := TestUser{
		Email:       "test@example.com",
		Password:    "Password123!",
		DisplayName: "Test User",
		Age:         25,
	}

	err := suite.validator.Validate(&user)
	assert.NoError(suite.T(), err)
}

func (suite *ValidatorTestSuite) TestValidateRequiredFields() {
	user := TestUser{
		Email:       "",
		Password:    "",
		DisplayName: "",
		Age:         25,
	}

	err := suite.validator.Validate(&user)
	assert.Error(suite.T(), err)

	ve, ok := err.(ValidationErrors)
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), ve.Errors, 3) // email, password, display_name required

	expectedFields := map[string]string{
		"email":           "email is required",
		testPasswordFieldConst: "password is required",
		"display_name":    "display_name is required",
	}

	for _, validationErr := range ve.Errors {
		expectedMsg, exists := expectedFields[validationErr.Field]
		assert.True(suite.T(), exists, "Unexpected field: %s", validationErr.Field)
		assert.Equal(suite.T(), expectedMsg, validationErr.Message)
		assert.Equal(suite.T(), "required", validationErr.Tag)
	}
}

func (suite *ValidatorTestSuite) TestValidateEmailFormat() {
	testCases := []struct {
		name  string
		email string
		valid bool
	}{
		{"Valid email", "test@example.com", true},
		{"Valid email with subdomain", "user@mail.example.com", true},
		{"Invalid email without @", "testexample.com", false},
		{"Invalid email without domain", "test@", false},
		{"Invalid email without local part", "@example.com", false},
		{"Empty email", "", false}, // This will fail on required, not email format
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			user := TestUser{
				Email:       tc.email,
				Password:    "Password123!",
				DisplayName: "Test User",
				Age:         25,
			}

			err := suite.validator.Validate(&user)

			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				ve, ok := err.(ValidationErrors)
				assert.True(t, ok)

				// Find email validation error
				hasEmailError := false
				for _, validationErr := range ve.Errors {
					if validationErr.Field == "email" && (validationErr.Tag == "email" || validationErr.Tag == "required") {
						hasEmailError = true
						break
					}
				}
				assert.True(t, hasEmailError, "Should have email validation error")
			}
		})
	}
}

func (suite *ValidatorTestSuite) TestValidateStringLength() {
	// Test minimum length
	user := TestUser{
		Email:       "test@example.com",
		Password:    "Short1!", // Too short (7 chars, min 8)
		DisplayName: "Test User",
		Age:         25,
	}

	err := suite.validator.Validate(&user)
	assert.Error(suite.T(), err)

	ve, ok := err.(ValidationErrors)
	assert.True(suite.T(), ok)

	hasMinLengthError := false
	for _, validationErr := range ve.Errors {
		if validationErr.Field == testPasswordFieldConst && validationErr.Tag == "min" {
			hasMinLengthError = true
			assert.Contains(suite.T(), validationErr.Message, "must be at least 8 characters long")
			break
		}
	}
	assert.True(suite.T(), hasMinLengthError)

	// Test maximum length - create a password that's too long but still complex
	longPassword := "A1!" + string(make([]byte, 126)) // 129 characters total (too long, max 128)
	for i := 3; i < len(longPassword); i++ {
		longPassword = longPassword[:i] + "a" + longPassword[i+1:]
	}

	user.Password = longPassword
	err = suite.validator.Validate(&user)
	assert.Error(suite.T(), err)
}

func (suite *ValidatorTestSuite) TestValidateNumericRange() {
	// Test minimum value
	user := TestUser{
		Email:       "test@example.com",
		Password:    "Password123!",
		DisplayName: "Test User",
		Age:         -1, // Invalid age
	}

	err := suite.validator.Validate(&user)
	assert.Error(suite.T(), err)

	ve, ok := err.(ValidationErrors)
	assert.True(suite.T(), ok)

	hasGteError := false
	for _, validationErr := range ve.Errors {
		if validationErr.Field == "age" && validationErr.Tag == "gte" {
			hasGteError = true
			assert.Contains(suite.T(), validationErr.Message, "must be greater than or equal to 0")
			break
		}
	}
	assert.True(suite.T(), hasGteError)

	// Test maximum value
	user.Age = 151 // Too old
	err = suite.validator.Validate(&user)
	assert.Error(suite.T(), err)

	ve, ok = err.(ValidationErrors)
	assert.True(suite.T(), ok)

	hasLteError := false
	for _, validationErr := range ve.Errors {
		if validationErr.Field == "age" && validationErr.Tag == "lte" {
			hasLteError = true
			assert.Contains(suite.T(), validationErr.Message, "must be less than or equal to 150")
			break
		}
	}
	assert.True(suite.T(), hasLteError)
}

func (suite *ValidatorTestSuite) TestValidateOneOf() {
	product := TestProduct{
		Name:        "Test Product",
		Price:       100,
		Category:    "invalid_category", // Should be one of: electronics, clothing, books
		Description: "Test description",
		Code:        "ABCD1234",
	}

	err := suite.validator.Validate(&product)
	assert.Error(suite.T(), err)

	ve, ok := err.(ValidationErrors)
	assert.True(suite.T(), ok)

	hasOneOfError := false
	for _, validationErr := range ve.Errors {
		if validationErr.Field == "category" && validationErr.Tag == "oneof" {
			hasOneOfError = true
			assert.Contains(suite.T(), validationErr.Message, "must be one of: electronics clothing books")
			break
		}
	}
	assert.True(suite.T(), hasOneOfError)
}

func (suite *ValidatorTestSuite) TestValidateAlphaNum() {
	product := TestProduct{
		Name:        "Test Product",
		Price:       100,
		Category:    "electronics",
		Description: "Test description",
		Code:        "ABC-123!", // Should be alphanumeric only
	}

	err := suite.validator.Validate(&product)
	assert.Error(suite.T(), err)

	ve, ok := err.(ValidationErrors)
	assert.True(suite.T(), ok)

	hasAlphaNumError := false
	for _, validationErr := range ve.Errors {
		if validationErr.Field == "code" && validationErr.Tag == "alphanum" {
			hasAlphaNumError = true
			assert.Contains(suite.T(), validationErr.Message, "must contain only alphanumeric characters")
			break
		}
	}
	assert.True(suite.T(), hasAlphaNumError)
}

func (suite *ValidatorTestSuite) TestValidateExactLength() {
	product := TestProduct{
		Name:        "Test Product",
		Price:       100,
		Category:    "electronics",
		Description: "Test description",
		Code:        "ABC123", // Should be exactly 8 characters
	}

	err := suite.validator.Validate(&product)
	assert.Error(suite.T(), err)

	ve, ok := err.(ValidationErrors)
	assert.True(suite.T(), ok)

	hasLenError := false
	for _, validationErr := range ve.Errors {
		if validationErr.Field == "code" && validationErr.Tag == "len" {
			hasLenError = true
			assert.Contains(suite.T(), validationErr.Message, "must be exactly 8 characters long")
			break
		}
	}
	assert.True(suite.T(), hasLenError)
}

func (suite *ValidatorTestSuite) TestValidationErrorsMessage() {
	user := TestUser{
		Email:       "invalid-email",
		Password:    "short",
		DisplayName: "",
		Age:         -1,
	}

	err := suite.validator.Validate(&user)
	assert.Error(suite.T(), err)

	ve, ok := err.(ValidationErrors)
	assert.True(suite.T(), ok)

	// Test that Error() method joins messages
	errorMsg := ve.Error()
	assert.Contains(suite.T(), errorMsg, ";")
	assert.NotEmpty(suite.T(), errorMsg)
}

func (suite *ValidatorTestSuite) TestJSONFieldNames() {
	// Test that field names come from JSON tags
	user := TestUser{
		DisplayName: "", // This should show as "display_name" not "DisplayName"
	}

	err := suite.validator.Validate(&user)
	assert.Error(suite.T(), err)

	ve, ok := err.(ValidationErrors)
	assert.True(suite.T(), ok)

	for _, validationErr := range ve.Errors {
		// Field names should come from JSON tags
		assert.NotContains(suite.T(), validationErr.Field, "DisplayName")
		if validationErr.Tag == "required" && validationErr.Value == "" {
			// Should use JSON field name
			assert.True(suite.T(),
				validationErr.Field == "display_name" ||
					validationErr.Field == "email" ||
					validationErr.Field == testPasswordFieldConst)
		}
	}
}

func TestValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}

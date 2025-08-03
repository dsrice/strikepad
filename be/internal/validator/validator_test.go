package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const testPasswordFieldConst = "password"

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
	Password    string `json:"password" validate:"required,min=8,max=128,password_complex"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
	Age         int    `json:"age" validate:"gte=0,lte=150"`
}

type TestProduct struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Category    string `json:"category" validate:"required,oneof=electronics clothing books"`
	Description string `json:"description" validate:"max=500"`
	Code        string `json:"code" validate:"required,alphanum,len=8"`
	Price       int    `json:"price" validate:"required,gt=0"`
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

				// For empty password, it will fail on RequiredTag first, not "password_complex"
				if tc.password == "" {
					hasRequiredError := false
					for _, validationErr := range ve.Errors {
						if validationErr.Tag == RequiredTag && validationErr.Field == testPasswordFieldConst {
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

func (suite *ValidatorTestSuite) TestValidateUser() {
	testCases := []struct {
		name        string
		user        TestUser
		expectError bool
		errorFields []string
	}{
		{
			name: "valid user",
			user: TestUser{
				Email:       "test@example.com",
				Password:    "Password123!",
				DisplayName: "Test User",
				Age:         25,
			},
			expectError: false,
		},
		{
			name: "all required fields missing",
			user: TestUser{
				Age: 25,
			},
			expectError: true,
			errorFields: []string{"email", "password", "display_name"},
		},
		{
			name: "invalid email format",
			user: TestUser{
				Email:       "invalid-email",
				Password:    "Password123!",
				DisplayName: "Test User",
				Age:         25,
			},
			expectError: true,
			errorFields: []string{"email"},
		},
		{
			name: "password too short",
			user: TestUser{
				Email:       "test@example.com",
				Password:    "Short1!",
				DisplayName: "Test User",
				Age:         25,
			},
			expectError: true,
			errorFields: []string{"password"},
		},
		{
			name: "age too young",
			user: TestUser{
				Email:       "test@example.com",
				Password:    "Password123!",
				DisplayName: "Test User",
				Age:         -1,
			},
			expectError: true,
			errorFields: []string{"age"},
		},
		{
			name: "age too old",
			user: TestUser{
				Email:       "test@example.com",
				Password:    "Password123!",
				DisplayName: "Test User",
				Age:         151,
			},
			expectError: true,
			errorFields: []string{"age"},
		},
		{
			name: "multiple validation errors",
			user: TestUser{
				Email:       "invalid-email",
				Password:    "short",
				DisplayName: "",
				Age:         -1,
			},
			expectError: true,
			errorFields: []string{"email", "password", "display_name", "age"},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.validator.Validate(&tc.user)

			if tc.expectError {
				assert.Error(t, err)
				ve, ok := err.(ValidationErrors)
				assert.True(t, ok)

				// Check that all expected fields have errors
				errorFieldsFound := make(map[string]bool)
				for _, validationErr := range ve.Errors {
					errorFieldsFound[validationErr.Field] = true
				}

				for _, expectedField := range tc.errorFields {
					assert.True(t, errorFieldsFound[expectedField], "Expected error for field: %s", expectedField)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *ValidatorTestSuite) TestValidateProduct() {
	testCases := []struct {
		name        string
		product     TestProduct
		expectError bool
		errorFields []string
	}{
		{
			name: "valid product",
			product: TestProduct{
				Name:        "Test Product",
				Price:       100,
				Category:    "electronics",
				Description: "Test description",
				Code:        "ABCD1234",
			},
			expectError: false,
		},
		{
			name: "invalid category (oneof)",
			product: TestProduct{
				Name:        "Test Product",
				Price:       100,
				Category:    "invalid_category",
				Description: "Test description",
				Code:        "ABCD1234",
			},
			expectError: true,
			errorFields: []string{"category"},
		},
		{
			name: "invalid code format (alphanum)",
			product: TestProduct{
				Name:        "Test Product",
				Price:       100,
				Category:    "electronics",
				Description: "Test description",
				Code:        "ABC-123!",
			},
			expectError: true,
			errorFields: []string{"code"},
		},
		{
			name: "invalid code length",
			product: TestProduct{
				Name:        "Test Product",
				Price:       100,
				Category:    "electronics",
				Description: "Test description",
				Code:        "ABC123",
			},
			expectError: true,
			errorFields: []string{"code"},
		},
		{
			name: "missing required fields",
			product: TestProduct{
				Description: "Test description",
			},
			expectError: true,
			errorFields: []string{"name", "category", "code", "price"},
		},
		{
			name: "invalid price (gt)",
			product: TestProduct{
				Name:        "Test Product",
				Price:       0,
				Category:    "electronics",
				Description: "Test description",
				Code:        "ABCD1234",
			},
			expectError: true,
			errorFields: []string{"price"},
		},
		{
			name: "description too long",
			product: TestProduct{
				Name:        "Test Product",
				Price:       100,
				Category:    "electronics",
				Description: strings.Repeat("a", 501),
				Code:        "ABCD1234",
			},
			expectError: true,
			errorFields: []string{"description"},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.validator.Validate(&tc.product)

			if tc.expectError {
				assert.Error(t, err)
				ve, ok := err.(ValidationErrors)
				assert.True(t, ok)

				// Check that all expected fields have errors
				errorFieldsFound := make(map[string]bool)
				for _, validationErr := range ve.Errors {
					errorFieldsFound[validationErr.Field] = true
				}

				for _, expectedField := range tc.errorFields {
					assert.True(t, errorFieldsFound[expectedField], "Expected error for field: %s", expectedField)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *ValidatorTestSuite) TestValidationErrorsAndFieldNames() {
	testCases := []struct {
		name           string
		user           TestUser
		expectedFields []string
		checkErrorMsg  bool
		checkJSONNames bool
	}{
		{
			name: "multiple validation errors",
			user: TestUser{
				Email:       "invalid-email",
				Password:    "short",
				DisplayName: "",
				Age:         -1,
			},
			expectedFields: []string{"email", "password", "display_name", "age"},
			checkErrorMsg:  true,
			checkJSONNames: true,
		},
		{
			name: "json field names test",
			user: TestUser{
				DisplayName: "",
			},
			expectedFields: []string{"email", "password", "display_name"},
			checkErrorMsg:  false,
			checkJSONNames: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.validator.Validate(&tc.user)
			assert.Error(t, err)

			ve, ok := err.(ValidationErrors)
			assert.True(t, ok)

			if tc.checkErrorMsg {
				// Test that Error() method joins messages
				errorMsg := ve.Error()
				assert.Contains(t, errorMsg, ";")
				assert.NotEmpty(t, errorMsg)
			}

			if tc.checkJSONNames {
				// Test that field names come from JSON tags
				for _, validationErr := range ve.Errors {
					// Field names should come from JSON tags, not struct field names
					assert.NotContains(t, validationErr.Field, "DisplayName")
					if validationErr.Tag == RequiredTag && validationErr.Value == "" {
						// Should use JSON field name
						assert.True(t,
							validationErr.Field == "display_name" ||
								validationErr.Field == "email" ||
								validationErr.Field == testPasswordFieldConst)
					}
				}
			}

			// Check that all expected fields have errors
			errorFieldsFound := make(map[string]bool)
			for _, validationErr := range ve.Errors {
				errorFieldsFound[validationErr.Field] = true
			}

			for _, expectedField := range tc.expectedFields {
				assert.True(t, errorFieldsFound[expectedField], "Expected error for field: %s", expectedField)
			}
		})
	}
}

func TestValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}
package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	RequiredTag = "required"
	EmailTag    = "email"
	PasswordTag = "password"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve ValidationErrors) Error() string {
	messages := make([]string, 0, len(ve.Errors))
	for _, err := range ve.Errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// Validator wraps the go-playground validator
type Validator struct {
	validator *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Register field name function to use JSON tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom password validation
	if err := v.RegisterValidation("password_complex", validatePasswordComplexity); err != nil {
		panic("Failed to register password_complex validation: " + err.Error())
	}

	return &Validator{
		validator: v,
	}
}

// validatePasswordComplexity validates that password contains lowercase, uppercase, and symbol
func validatePasswordComplexity(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)

	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)

	// Check for at least one symbol (non-alphanumeric character)
	hasSymbol := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return hasLower && hasUpper && hasSymbol
}

// Validate validates a struct and returns formatted errors
func (v *Validator) Validate(s interface{}) error {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		// If it's not a ValidationErrors type, return a generic validation error
		return ValidationErrors{
			Errors: []ValidationError{{
				Field:   "unknown",
				Tag:     "validation",
				Value:   "",
				Message: err.Error(),
			}},
		}
	}
	validationErrors := make([]ValidationError, 0, len(validationErrs))

	for _, err := range validationErrs {
		ve := ValidationError{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   fmt.Sprintf("%v", err.Value()),
			Message: getErrorMessage(err),
		}
		validationErrors = append(validationErrors, ve)
	}

	return ValidationErrors{Errors: validationErrors}
}

// getErrorMessage returns a human-readable error message for validation errors
func getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()

	if msg := getBasicValidationMessage(field, tag); msg != "" {
		return msg
	}

	if msg := getLengthValidationMessage(field, tag, fe); msg != "" {
		return msg
	}

	if msg := getNumericValidationMessage(field, tag, fe.Param()); msg != "" {
		return msg
	}

	if msg := getFormatValidationMessage(field, tag, fe.Param()); msg != "" {
		return msg
	}

	return fmt.Sprintf("%s is invalid", field)
}

// getBasicValidationMessage handles basic validation messages
func getBasicValidationMessage(field, tag string) string {
	switch tag {
	case RequiredTag:
		return fmt.Sprintf("%s is required", field)
	case EmailTag:
		return fmt.Sprintf("%s must be a valid email address", field)
	case "password_complex":
		return fmt.Sprintf("%s must contain at least one lowercase letter, one uppercase letter, and one symbol", field)
	default:
		return ""
	}
}

// getLengthValidationMessage handles length-related validation messages
func getLengthValidationMessage(field, tag string, fe validator.FieldError) string {
	switch tag {
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
		}
		return fmt.Sprintf("%s must be at least %s", field, fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters long", field, fe.Param())
		}
		return fmt.Sprintf("%s must be at most %s", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, fe.Param())
	default:
		return ""
	}
}

// getNumericValidationMessage handles numeric comparison validation messages
func getNumericValidationMessage(field, tag, param string) string {
	switch tag {
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	default:
		return ""
	}
}

// getFormatValidationMessage handles format and content validation messages
func getFormatValidationMessage(field, tag, param string) string {
	switch tag {
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be a number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", field)
	default:
		return ""
	}
}

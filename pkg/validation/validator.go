package validation

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	// Register custom validations
	v.RegisterValidation("date", validateDate)
	v.RegisterValidation("phone", validatePhone)

	return &Validator{validate: v}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

func validateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return true // empty is valid (optional field)
	}

	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // empty is valid (optional field)
	}

	// Basic phone validation - adjust regex as needed
	phoneRegex := `^\+?[1-9]\d{1,14}$`
	matched, _ := regexp.MatchString(phoneRegex, phone)
	return matched
}

func (v *Validator) FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()

			switch tag {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errors[field] = fmt.Sprintf("%s must be a valid email", field)
			case "date":
				errors[field] = fmt.Sprintf("%s must be a valid date in YYYY-MM-DD format", field)
			case "phone":
				errors[field] = fmt.Sprintf("%s must be a valid phone number", field)
			default:
				errors[field] = fmt.Sprintf("%s failed %s validation", field, tag)
			}
		}
	}

	return errors
}

package util

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func EmailValidator(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return validateEmail(email)
}
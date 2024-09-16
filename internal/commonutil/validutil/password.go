package validutil

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	atLeastTenChars := len(password) >= 10
	hasDigit := strings.ContainsFunc(password, unicode.IsDigit)
	hasChar := strings.ContainsFunc(password, unicode.IsLetter)

	return atLeastTenChars && hasDigit && hasChar
}

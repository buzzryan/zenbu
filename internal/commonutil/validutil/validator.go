package validutil

import (
	"log"

	"github.com/go-playground/validator/v10"
)

var val *validator.Validate

func init() {
	val = validator.New(validator.WithRequiredStructEnabled())
	if err := val.RegisterValidation("password", passwordValidator); err != nil {
		log.Panicf("failed to register password validator: %v", err)
	}
}

func Validate(v interface{}) error {
	return val.Struct(v)
}

package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

func ValidateStruct(model interface{}) error {
	err := validate.Struct(model)
	return err
}

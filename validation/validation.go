package validation

import (
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

func validatePhone(fieldLevel validator.FieldLevel) bool {
	result, _ := regexp.MatchString(`^09\d{5,9}$`, fieldLevel.Field().String())

	return result
}

func validateAlphaSpace(fieldLevel validator.FieldLevel) bool {
	result, _ := regexp.MatchString(`^[\pL\pM][\pL\pM\s‌]+[\pL\pM]$`, fieldLevel.Field().String())

	return result
}

func GetValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("alpha_space", validateAlphaSpace)
	return validate
}

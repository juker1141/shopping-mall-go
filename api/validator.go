package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validFullName validator.Func = func(fieldLevel validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[a-zA-Z\s\p{Han}]+$`)
	return re.MatchString(fieldLevel.Field().String())
}

var validCellPhone validator.Func = func(fieldLevel validator.FieldLevel) bool {
	phone := fieldLevel.Field().String()

	pattern := `^(09\d{8}|\+8869\d{8})$`
	matched, _ := regexp.MatchString(pattern, phone)

	return matched
}

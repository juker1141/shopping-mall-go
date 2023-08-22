package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validFullName validator.Func = func(fieldLevel validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[a-zA-Z\s\p{Han}]+$`)
	return re.MatchString(fieldLevel.Field().String())
}

package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validFullName validator.Func = func(fieldLevel validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[a-zA-Z\s\p{Han}]+$`)
	return re.MatchString(fieldLevel.Field().String())
}

var validTaiwanPhone validator.Func = func(fieldLevel validator.FieldLevel) bool {
	phone := fieldLevel.Field().String()

	// 使用正则表达式验证台湾电话号码格式
	// 此正则表达式匹配台湾的手机格式，例如：09xx-xxxxxx 或 +8869xx-xxxxxx
	pattern := `^(09\d{2}-?\d{6}|(\+8869\d{2}-?)?\d{6}|0[1-9]\d{7})$`
	matched, _ := regexp.MatchString(pattern, phone)

	return matched
}

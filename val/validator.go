package val

import (
	"fmt"
	"regexp"
)

var (
	isValidAccount  = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s\p{Han}]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateAccount(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidAccount(value) {
		return fmt.Errorf("must contain only lowercase letters, digits")
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("must only contain Chinese characters, English letters, or spaces")
	}
	return nil
}

func IsValidStatus(status int) bool {
	switch status {
	case 0:
		return true
	case 1:
		return true
	default:
		return false
	}
}

func ValidateStatus(status int) error {
	if !IsValidStatus(status) {
		return fmt.Errorf("must be 0 or 1")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

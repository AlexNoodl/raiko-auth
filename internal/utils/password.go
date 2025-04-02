package utils

import (
	"regexp"
)

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	hasUpper := regexp.MustCompile("[A-Z]").MatchString(password)
	hasLower := regexp.MustCompile("[a-z]").MatchString(password)
	hasNumber := regexp.MustCompile("[0-9]").MatchString(password)
	hasSpecial := regexp.MustCompile("[!@#$%^&*()]").MatchString(password)

	return hasLower && hasUpper && hasNumber && hasSpecial
}

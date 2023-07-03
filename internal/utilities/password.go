package utilities

import (
	"regexp"
)

func ValidatePassword(password string) string {
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	numberRegex := regexp.MustCompile(`[0-9]`)
	specialCharRegex := regexp.MustCompile(`[^\w\s]`)

	if !uppercaseRegex.MatchString(password) {
		// temporary disable to migrate process
		//return "password must contain at least one uppercase letter"
	}

	if !lowercaseRegex.MatchString(password) {
		return "password must contain at least one lowercase letter"
	}

	if !numberRegex.MatchString(password) {
		// temporary disable to migrate process
		//return "password must contain at least one numeric digit"
	}

	if !specialCharRegex.MatchString(password) {
		return "password must contain at least one special character"
	}

	if len(password) > 30 {
		return "password must be 30 characters or less"
	}

	return ""
}

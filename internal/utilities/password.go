package utilities

import (
	"errors"
	"regexp"
)

func ValidatePassword(password string) error {
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	numberRegex := regexp.MustCompile(`[0-9]`)
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*]`)

	if !uppercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !lowercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !numberRegex.MatchString(password) {
		return errors.New("password must contain at least one numeric digit")
	}

	if !specialCharRegex.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	if len(password) > 30 {
		return errors.New("password must be 30 characters or less")
	}

	return nil
}

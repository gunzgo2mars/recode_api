package validation

import (
	"regexp"
)

func ValidateEmail(email string) bool {

	return regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email)

}

func ValidateLetter(text string) bool {

	return regexp.MustCompile(`^[a-zA-Zก-ฮ]*$`).MatchString(text)

}

/// Ignition Start !!!!

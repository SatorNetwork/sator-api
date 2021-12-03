package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mcnijman/go-emailaddress"
)

var (
	ErrInvalidEmailAddress = errors.New("invalid email address")
	ErrInvalidIcanSuffix   = errors.New("invalid ICAN suffix")
)

// SanitizeEmail cleans email address from dots, dashes, etc
func SanitizeEmail(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	email, err := emailaddress.Parse(s)
	if err != nil {
		return "", ErrInvalidEmailAddress
	}

	username := email.LocalPart
	domain := email.Domain

	if strings.Contains(username, "+") {
		p := strings.Split(username, "+")
		username = p[0]
	}

	if strings.Contains(username, ".") {
		p := strings.Split(username, ".")
		username = strings.Join(p, "")
	}

	if strings.Contains(username, "-") {
		p := strings.Split(username, "-")
		username = strings.Join(p, "")
	}

	result := fmt.Sprintf("%s@%s", username, domain)

	if err := email.ValidateIcanSuffix(); err != nil {
		return result, ErrInvalidIcanSuffix
	}

	return result, nil
}

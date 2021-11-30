package utils

import (
	"fmt"
	"strings"
)

// SanitizeEmail cleans email address from dots, dashes, etc
func SanitizeEmail(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email address")
	}

	username := parts[0]
	domain := parts[1]

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

	return fmt.Sprintf("%s@%s", username, domain), nil
}

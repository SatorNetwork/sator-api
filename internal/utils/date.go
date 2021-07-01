package utils

import (
	"fmt"
	"time"
)

// DateFromString is used to parse date from string according to layout.
func DateFromString(date string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse date from string:%w", err)
	}

	return t, nil
}

package utils

import (
	"fmt"
	"time"
)

// DateLayout sets the format for the date.
const DateLayout = "2006-01-02T15:04:05.000Z"

// DateFromString is used to parse date from string according to layout.
func DateFromString(date string) (time.Time, error) {
	t, err := time.Parse(DateLayout, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse date from string:%w", err)
	}

	return t, nil
}

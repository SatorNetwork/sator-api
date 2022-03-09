package validator

import (
	"github.com/mcnijman/go-emailaddress"
)

func ValidateEmail(s string) error {
	email, err := emailaddress.Parse(s)
	if err != nil {
		return err
	}

	// if err := email.ValidateHost(); err != nil {
	// 	return err
	// }

	if err := email.ValidateIcanSuffix(); err != nil {
		return err
	}

	return nil
}

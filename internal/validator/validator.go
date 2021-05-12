package validator

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	// ValidationError struct
	ValidationError struct {
		ErrorsBag url.Values
		ErrorStr  string
	}

	// ValidateFunc interface
	ValidateFunc func(s interface{}) error
)

func (v ValidationError) Error() string {
	return v.ErrorStr
}

// ValidateStruct validates structures with rules in validate tags,
// return nil or instance of the ValidationError structure
// as validator it uses package: github.com/go-playground/validator/v10
func ValidateStruct() func(s interface{}) error {
	// Init validator
	var (
		errStr   string
		validate = validator.New()
	)

	return func(s interface{}) error {
		errs := url.Values{}
		err := validate.Struct(s)
		if err != nil {
			switch e := err.(type) {
			case *validator.InvalidValidationError:
				return fmt.Errorf("invalid validation error: %w", e)
			case validator.ValidationErrors:
				errStr = err.Error()
				for _, verr := range e {
					errs.Add(strings.ToLower(verr.Field()), verr.Tag())
				}
			}

			if len(errs) > 0 {
				return ValidationError{
					ErrorsBag: errs,
					ErrorStr:  errStr,
				}
			}

			return fmt.Errorf("internal validation error: %w", err)
		}
		return nil
	}
}

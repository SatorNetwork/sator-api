package utils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/thedevsaddam/govalidator"
)

// ValidationError ...
type ValidationError struct {
	errBag url.Values
}

// NewValidationError factory
func NewValidationError(errs url.Values) *ValidationError {
	return &ValidationError{errBag: errs}
}

// Implementation of the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validationError: %v", e.errBag)
}

// Validate validation helper, returns HTTP error or nil
func Validate(req *http.Request, rul, msg map[string][]string) error {
	v := govalidator.New(govalidator.Options{
		Request:  req,
		Rules:    rul,
		Messages: msg,
	})
	if err := v.Validate(); len(err) > 0 {
		return NewValidationError(err)
	}

	return nil
}

package handler

import (
	"net/http"

	"github.com/thedevsaddam/govalidator"
)

// Validation helper,
// returns HTTP error or nil
func validate(req *http.Request, rul, msg map[string][]string) error {
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

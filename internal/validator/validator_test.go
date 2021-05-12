package validator

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	vf := ValidateStruct()

	s := struct {
		Test string `json:"test" validate:"required"`
	}{}

	if err := vf(&s); err == nil {
		t.Errorf("Expected 'required' validation error, got=%v", err)
	}

	s.Test = "test val"
	if err := vf(&s); err != nil {
		t.Errorf("Expected 'required' validation error, got=%v", err)
	}
}

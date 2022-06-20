package alias

import "github.com/pkg/errors"

type Value string

const (
	FlagValueUndefined Value = "UNDEFINED"
	FlagValueEnabled   Value = "ENABLED"
	FlagValueDisabled  Value = "DISABLED"
)

func NewFlagValueFromString(s string) (Value, error) {
	switch s {
	case "ENABLED":
		return FlagValueEnabled, nil
	case "DISABLED":
		return FlagValueDisabled, nil
	default:
		return FlagValueUndefined, errors.Errorf("flags value with such name %v doesn't exist", s)
	}
}

func (v Value) String() string {
	switch v {
	case FlagValueEnabled:
		return "ENABLED"
	case FlagValueDisabled:
		return "DISABLED"
	default:
		return "UNDEFINED"
	}
}

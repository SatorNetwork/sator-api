package alias

import "github.com/pkg/errors"

type Alias uint8

const (
	UndefinedAlias Alias = iota
	FeePayerAlias
	TokenHolderAlias
)

func newAliasFromString(s string) (Alias, error) {
	switch s {
	case "fee_payer":
		return FeePayerAlias, nil
	case "token_holder":
		return TokenHolderAlias, nil
	default:
		return UndefinedAlias, errors.Errorf("types with such name %v doesn't exist", s)
	}
}

func (a Alias) String() string {
	switch a {
	case UndefinedAlias:
		return "undefined"
	case FeePayerAlias:
		return "fee_payer"
	case TokenHolderAlias:
		return "token_holder"
	default:
		return "undefined"
	}
}

type Aliases []Alias

func NewAliasesFromStrings(strings []string) (Aliases, error) {
	aliases := make(Aliases, 0, len(strings))
	for _, s := range strings {
		alias, err := newAliasFromString(s)
		if err != nil {
			return nil, errors.Wrap(err, "can't get new types from string")
		}
		aliases = append(aliases, alias)
	}

	return aliases, nil
}

func (as Aliases) ToStrings() []string {
	strings := make([]string, 0, len(as))
	for _, a := range as {
		strings = append(strings, a.String())
	}

	return strings
}

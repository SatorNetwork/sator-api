package gapi

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/golang-jwt/jwt"
)

// SignStruct function is used to sign a struct
// using the given signing key.
// Returns a signature string or an error.
func SignStruct(signingKey []byte, s interface{}) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("could not marshal struct: %w", err)
	}

	sig, err := jwt.SigningMethodHS256.Sign(string(b), signingKey)
	if err != nil {
		return "", ErrCouldNotSignResponse
	}

	return sig, nil
}

// ValidateSignature function is used to validate a signature
// using the given signing key.
// Returns an error if the signature is invalid.
func ValidateSignature(signingKey []byte, s interface{}, sig string) error {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("could not marshal struct: %v", err)
	}

	if err := jwt.SigningMethodHS256.Verify(string(b), sig, signingKey); err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return err
		}

		return fmt.Errorf("%w: %s", ErrCouldNotVerifySignature, err.Error())
	}

	return nil
}

// SignResponse function is used to sign a response
// using the given signing key.
// Returns a signature string or an error.
func SignResponse(signingKey []byte, response interface{}) (string, error) {
	if response == nil {
		return "", nil
	}

	var structToSign interface{}

	switch r := response.(type) {
	case httpencoder.Response, httpencoder.BoolResultResponse, httpencoder.ListResponse:
		structToSign = r
	case bool:
		structToSign = httpencoder.BoolResult(r)
	default:
		structToSign = httpencoder.Response{Data: response}
	}

	return SignStruct(signingKey, structToSign)
}

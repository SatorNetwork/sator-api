package gapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
)

// SignStruct function is used to sign a struct
// using the given signing key.
// Returns a signature string or an error.
func SignStruct(signingKey []byte, s interface{}) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", ErrCouldNotSignResponse
	}

	return signHS256(b, signingKey), nil
}

// ValidateSignature function is used to validate a signature
// using the given signing key.
// Returns an error if the signature is invalid.
func ValidateSignature(signingKey []byte, s interface{}, sig string) error {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("could not marshal struct: %v", err)
	}

	verified, err := verifyHS256(b, signingKey, sig)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCouldNotVerifySignature, err.Error())
	}

	if !verified {
		return ErrSignatureInvalid
	}

	return nil
}

// ValidateStringSignature function is used to validate a signature
// using the given signing key.
// Returns an error if the signature is invalid.
func ValidateStringSignature(signingKey []byte, s string, sig string) error {
	verified, err := verifyHS256([]byte(s), signingKey, sig)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCouldNotVerifySignature, err.Error())
	}

	if !verified {
		return ErrSignatureInvalid
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

// signHS256 signs a message using the HS256 algorithm.
func signHS256(msg, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)

	return hex.EncodeToString(mac.Sum(nil))
}

// verifyHS256 verifies a message using the HS256 algorithm.
func verifyHS256(msg, key []byte, hash string) (bool, error) {
	sig, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(msg)

	return hmac.Equal(sig, mac.Sum(nil)), nil
}

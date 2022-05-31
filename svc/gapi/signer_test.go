package gapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignAndValidateStruct(t *testing.T) {
	type TestStruct struct {
		A string `json:"a"`
		B int    `json:"b"`
	}

	testStruct := TestStruct{
		A: "test",
		B: 1,
	}

	signature, err := SignStruct([]byte("secret"), testStruct)
	assert.NoError(t, err)
	assert.NotEmpty(t, signature)
	t.Logf("signature: %s", signature)

	err = ValidateSignature([]byte("secret"), testStruct, signature)
	assert.NoError(t, err)
}

// TestSignAndValidateString tests signing and validating a string.
func TestSignAndValidateString(t *testing.T) {
	payload := []byte(`{"blocks_done":1,"game_result":0}`)
	// payload := []byte("123")
	signature := signHS256(payload, []byte("secret"))
	assert.NotEmpty(t, signature)
	t.Logf("signature: %s", signature)

	valid, err := verifyHS256(payload, []byte("secret"), signature)
	assert.NoError(t, err)
	assert.True(t, valid)
}

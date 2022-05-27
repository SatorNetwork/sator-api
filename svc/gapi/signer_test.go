package gapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSignature(t *testing.T) {
	unityVersion := "2020.3.14f1"
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMThkNjdkMGItOWJlOS00NzU5LWI2NzktZGJiNTI3M2JhOWFjIiwidXNlcm5hbWUiOiJxNDgwODQwMiIsInJvbGUiOiJhZG1pbiIsImRldmljZV9pZCI6ImNtQVRLU18zU3EyNlJuRGY0NWtfWW0iLCJleHAiOjE2NTM3NjE3MjAsImp0aSI6IjdlMmY3N2I1LTRjMDctNGUxNy04ZGMyLTIyMTQ3ZGQzOGUzYyIsImlhdCI6MTY1MzY3NTMyMCwibmJmIjoxNjUzNjc1MzIwLCJzdWIiOiJhY2Nlc3NfdG9rZW4ifQ.UiXfSOIYhT0SCVJneUQ76ngKQNOUPaqbY7sUfVNWxeM"
	timestamp := "1653676279"
	payload := `{"blocks_done":1,"game_result":0}`
	deviceID := "cmATKS_3Sq26RnDf45k_Ym"
	signature := "S00CKC2mvx2TZNktTRN5p8mN3DB2-csWtaYVn47qeUc"

	signingKey := fmt.Sprintf("%s%s%s%s", unityVersion, deviceID, timestamp, jwt[len(jwt)-3:])
	_ = signingKey

	err := ValidateStringSignature([]byte("secret"), payload, signature)
	assert.NoError(t, err)
}

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

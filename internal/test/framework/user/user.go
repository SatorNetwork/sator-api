package user

import (
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"

	internal_rsa "github.com/SatorNetwork/sator-api/internal/encryption/rsa"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/auth"
)

type User struct {
	signUpRequest *auth.SignUpRequest

	email       string
	username    string
	password    string
	accessToken string

	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	c *client.Client
	t *testing.T
}

func NewInitializedUser(signUpRequest *auth.SignUpRequest, t *testing.T) *User {
	user := &User{
		signUpRequest: signUpRequest,

		email:    signUpRequest.Email,
		username: signUpRequest.Username,
		password: signUpRequest.Password,

		c: client.NewClient(),
		t: t,
	}

	user.SignUp()
	user.VerifyAccount()
	user.RegisterPublicKey()

	return user
}

func (u *User) Username() string {
	return u.username
}

func (u *User) AccessToken() string {
	return u.accessToken
}

func (u *User) PrivateKey() *rsa.PrivateKey {
	return u.privateKey
}

func (u *User) SignUp() {
	signUpResp, err := u.c.Auth.SignUp(u.signUpRequest)
	require.NoError(u.t, err)
	require.NotNil(u.t, signUpResp)
	require.NotEmpty(u.t, signUpResp.AccessToken)

	u.accessToken = signUpResp.AccessToken
}

func (u *User) VerifyAccount() {
	err := u.c.Auth.VerifyAcount(u.accessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(u.t, err)
}

func (u *User) RegisterPublicKey() {
	var err error
	u.privateKey, u.publicKey, err = internal_rsa.GenerateKeyPair(4096)
	require.NoError(u.t, err)
	publcKeyBytes, err := internal_rsa.PublicKeyToBytes(u.publicKey)
	require.NoError(u.t, err)
	err = u.c.Auth.RegisterPublicKey(u.accessToken, &auth.RegisterPublicKeyRequest{
		PublicKey: string(publcKeyBytes),
	})
	require.NoError(u.t, err)
}
package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
)

const (
	defaultEmail1    = "john.doe@sator.io"
	defaultPassword1 = "qwerty12345"
)

type AuthClient struct{}

func New() *AuthClient {
	return new(AuthClient)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	AccessToken string `json:"access_token"`
}

type VerifyAccountRequest struct {
	OTP string `json:"otp"`
}

type RegisterPublicKeyRequest struct {
	PublicKey string `json:"public_key"`
}

func RandomSignUpRequest() *SignUpRequest {
	rand.Seed(time.Now().UnixNano())
	n := rand.Uint64()
	email := fmt.Sprintf("john.doe%v@sator.io", n)
	username := fmt.Sprintf("johndoe%v", n)
	return &SignUpRequest{
		Email:    email,
		Username: username,
		Password: defaultPassword1,
	}
}

func (a *AuthClient) Login(req *LoginRequest) (*LoginResponse, error) {
	url := "http://localhost:8080/auth/login"
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal create transfer request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Device-ID", uuid.New().String())
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp LoginResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}

func (a *AuthClient) SignUp(req *SignUpRequest) (*SignUpResponse, error) {
	url := "http://localhost:8080/auth/signup"
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal create transfer request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Device-ID", uuid.New().String())
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp SignUpResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}

func (a *AuthClient) VerifyAcount(accessToken string, req *VerifyAccountRequest) error {
	url := "http://localhost:8080/auth/verify-account"
	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can't marshal create transfer request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "can't read response body")
	}
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	return nil
}

func (a *AuthClient) RegisterPublicKey(accessToken string, req *RegisterPublicKeyRequest) error {
	url := "http://localhost:8080/auth/user/public_key/register"
	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can't marshal register public key request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "can't read response body")
	}
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	return nil
}

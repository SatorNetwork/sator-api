package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	client_utils "github.com/SatorNetwork/sator-api/internal/test/framework/client/utils"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultEmail1    = "john.doe@mail.dev"
	defaultPassword1 = "qwerty12345"
)

type AuthClient struct{}

func NewAuthClient() *AuthClient {
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

func RandomSignUpRequest() *SignUpRequest {
	rand.Seed(time.Now().UnixNano())
	n := rand.Uint64()
	email := fmt.Sprintf("john.doe%v@mail.dev", n)
	username := fmt.Sprintf("johndoe%v", n)
	return &SignUpRequest{
		Email:    email,
		Username: username,
		Password: defaultPassword1,
	}
}

func (a *AuthClient) Login(req *LoginRequest) (*LoginResponse, error) {
	url := fmt.Sprintf("http://localhost:8080/auth/login")
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal create transfer request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
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
	url := fmt.Sprintf("http://localhost:8080/auth/signup")
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal create transfer request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
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
	url := fmt.Sprintf("http://localhost:8080/auth/verify-account")
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

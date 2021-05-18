package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Post("/login", httptransport.NewServer(
		e.Login,
		decodeLoginRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/logout", httptransport.NewServer(
		e.Logout,
		decodeLogoutRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/refresh-token", httptransport.NewServer(
		e.RefreshToken,
		decodeRefreshTokenRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/signup", httptransport.NewServer(
		e.SignUp,
		decodeSignUpRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/forgot-password", httptransport.NewServer(
		e.ForgotPassword,
		decodeForgotPasswordRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/validate-reset-password-code", httptransport.NewServer(
		e.ValidateResetPasswordCode,
		decodeValidateResetPasswordCodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/reset-password", httptransport.NewServer(
		e.ResetPassword,
		decodeResetPasswordRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/verify-account", httptransport.NewServer(
		e.VerifyAccount,
		decodeVerifyAccountRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeLogoutRequest(ctx context.Context, _ *http.Request) (request interface{}, err error) {
	tid, err := jwt.TokenIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user id: %w", err)
	}

	return tid.String(), nil
}

func decodeRefreshTokenRequest(ctx context.Context, _ *http.Request) (request interface{}, err error) {
	uid, err := jwt.UserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user id: %w", err)
	}

	tid, err := jwt.TokenIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user id: %w", err)
	}

	return RefreshTokenRequest{
		UserID:  uid.String(),
		TokenID: tid.String(),
	}, nil
}

func decodeSignUpRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeForgotPasswordRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeValidateResetPasswordCodeRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req ValidateResetPasswordCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeResetPasswordRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeVerifyAccountRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req VerifyAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	userID, err := jwt.UserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get user id: %w", err)
	}
	req.UserID = userID.String()

	return req, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidCredentials) {
		return http.StatusUnauthorized, err.Error()
	}

	if errors.Is(err, ErrEmailAlreadyTaken) ||
		errors.Is(err, ErrEmailAlreadyVerified) ||
		errors.Is(err, ErrOTPCode) {
		return http.StatusBadRequest, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

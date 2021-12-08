package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/deviceid"
	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/SatorNetwork/sator-api/internal/utils"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Predefined request query keys
const (
	allowedValue    = "allowed_value"
	restrictedValue = "restricted_value"
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
		httptransport.ServerBefore(jwtkit.HTTPToContext(), deviceid.ToContext()),
	}

	r.Get("/", httptransport.NewServer(
		e.Auth,
		decodeAuthRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/login", httptransport.NewServer(
		e.Login,
		decodeLoginRequest,
		encodeTokenResponse,
		options...,
	).ServeHTTP)

	r.Post("/logout", httptransport.NewServer(
		e.Logout,
		decodeLogoutRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/refresh-token", httptransport.NewServer(
		e.RefreshToken,
		decodeRefreshTokenRequest,
		encodeTokenResponse,
		options...,
	).ServeHTTP)

	r.Post("/refresh-token", httptransport.NewServer(
		e.RefreshToken,
		decodeRefreshTokenRequest,
		encodeTokenResponse,
		options...,
	).ServeHTTP)

	r.Post("/signup", httptransport.NewServer(
		e.SignUp,
		decodeSignUpRequest,
		encodeTokenResponse,
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

	r.Post("/change-password", httptransport.NewServer(
		e.ChangePassword,
		decodeChangePasswordRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/verify-account", httptransport.NewServer(
		e.VerifyAccount,
		decodeVerifyAccountRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/request-destroy-account", httptransport.NewServer(
		e.RequestDestroyAccount,
		decodeRequestDestroyAccountRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/verify-destroy-account-code", httptransport.NewServer(
		e.VerifyDestroyCode,
		decodeVerifyDestroyAccountRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/destroy", httptransport.NewServer(
		e.DestroyAccount,
		decodeDestroyAccountRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/is-verified", httptransport.NewServer(
		e.IsVerified,
		decodeIsVerifiedRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/resend-otp", httptransport.NewServer(
		e.ResendOTP,
		decodeResendOTPRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/request-update-email", httptransport.NewServer(
		e.RequestChangeEmail,
		decodeRequestChangeEmailRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/update-email", httptransport.NewServer(
		e.UpdateEmail,
		decodeUpdateEmailRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/update-username", httptransport.NewServer(
		e.UpdateUsername,
		decodeUpdateUsernameRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/whitelist", httptransport.NewServer(
		e.GetWhitelist,
		decodeGetWhitelist,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/whitelist", httptransport.NewServer(
		e.AddToWhitelist,
		decodeEditWhitelist,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/whitelist", httptransport.NewServer(
		e.DeleteFromWhitelist,
		decodeEditWhitelist,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/blacklist", httptransport.NewServer(
		e.GetBlacklist,
		decodeGetBlacklist,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/blacklist", httptransport.NewServer(
		e.AddToBlacklist,
		decodeEditBlacklist,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/blacklist", httptransport.NewServer(
		e.DeleteFromBlacklist,
		decodeEditBlacklist,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/kyc/access_token", httptransport.NewServer(
		e.GetAccessTokenByUserID,
		decodeGetAccessTokenByUserID,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/kyc/callback", httptransport.NewServer(
		e.VerificationCallback,
		decodeVerificationCallBack,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeAuthRequest(ctx context.Context, _ *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeLogoutRequest(ctx context.Context, _ *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeRefreshTokenRequest(ctx context.Context, _ *http.Request) (request interface{}, err error) {
	return nil, nil
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

func decodeChangePasswordRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}
func decodeVerifyAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	return req, nil
}

func decodeRequestChangeEmailRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req RequestChangeEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeUpdateEmailRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req UpdateEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeUpdateUsernameRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req UpdateUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeRequestDestroyAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeVerifyDestroyAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	return req, nil
}

func decodeDestroyAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	return req, nil
}

func decodeIsVerifiedRequest(ctx context.Context, _ *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeResendOTPRequest(ctx context.Context, _ *http.Request) (request interface{}, err error) {
	return nil, nil
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

func encodeTokenResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set(httpencoder.ContentTypeHeader, httpencoder.ContentType)
	return json.NewEncoder(w).Encode(response)
}

func decodeGetWhitelist(_ context.Context, r *http.Request) (interface{}, error) {
	return GetWhitelistRequest{
		AllowedValue: r.URL.Query().Get(allowedValue),
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
		},
	}, nil
}

func decodeEditWhitelist(_ context.Context, r *http.Request) (interface{}, error) {
	var req WhitelistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeGetBlacklist(_ context.Context, r *http.Request) (interface{}, error) {
	return GetBlacklistRequest{
		RestrictedValue: r.URL.Query().Get(restrictedValue),
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
		},
	}, nil
}

func decodeEditBlacklist(_ context.Context, r *http.Request) (interface{}, error) {
	var req BlacklistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeGetAccessTokenByUserID(ctx context.Context, _ *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeVerificationCallBack(_ context.Context, r *http.Request) (interface{}, error) {
	var req VerificationCallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req.ExternalUserId, nil
}

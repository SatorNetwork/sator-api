package auth

import (
	"context"

	"github.com/SatorNetwork/sator-api/internal/validator"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints struct
	Endpoints struct {
		Login          endpoint.Endpoint
		Logout         endpoint.Endpoint
		SignUp         endpoint.Endpoint
		ForgotPassword endpoint.Endpoint
		ResetPassword  endpoint.Endpoint
		VerifyAccount  endpoint.Endpoint
	}

	authService interface {
		Login(ctx context.Context, email, password string) (string, error)
		Logout(ctx context.Context) error
		SignUp(ctx context.Context, email, password, username string) error
		ForgotPassword(ctx context.Context, email string) error
		ResetPassword(ctx context.Context, email, password, otp string) error
		VerifyAccount(ctx context.Context) error
	}

	// AccessToken struct
	AccessToken struct {
		Token string `json:"access_token"`
	}

	// LoginRequest struct
	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// SignUpRequest struct
	SignUpRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		Username string `json:"user_n	ame" validate:"required"`
	}

	// ForgotPasswordRequest struct
	ForgotPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	// ResetPasswordRequest struct
	ResetPasswordRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		OTP      string `json:"otp" validate:"required"`
	}
)

// MakeEndpoints ...
func MakeEndpoints(as authService, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		Login:          MakeLoginEndpoint(as, validateFunc),
		Logout:         MakeLogoutEndpoint(as, validateFunc),
		SignUp:         MakeSignUpEndpoint(as, validateFunc),
		ForgotPassword: MakeForgotPasswordEndpoint(as, validateFunc),
		ResetPassword:  MakeResetPasswordEndpoint(as, validateFunc),
	}

	if len(m) > 0 {
		for _, mdw := range m {
			e.Login = mdw(e.Login)
		}
	}

	return e
}

// MakeLoginEndpoint ...
func MakeLoginEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		token, err := s.Login(ctx, req.Email, req.Password)
		if err != nil {
			return nil, err
		}

		return AccessToken{Token: token}, nil
	}
}

// MakeLogoutEndpoint ...
func MakeLogoutEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.Logout(ctx)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// MakeSignUpEndpoint ...
func MakeSignUpEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		err := s.SignUp(ctx, req.Email, req.Password, req.Username)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// MakeForgotPasswordEndpoint ...
func MakeForgotPasswordEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ForgotPasswordRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		err := s.ForgotPassword(ctx, req.Email)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// MakeResetPasswordEndpoint ...
func MakeResetPasswordEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ResetPasswordRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		err := s.ResetPassword(ctx, req.Email, req.Password, req.OTP)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

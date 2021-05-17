package auth

import (
	"context"

	"github.com/SatorNetwork/sator-api/internal/validator"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints struct
	Endpoints struct {
		Login                     endpoint.Endpoint
		Logout                    endpoint.Endpoint
		SignUp                    endpoint.Endpoint
		ForgotPassword            endpoint.Endpoint
		ValidateResetPasswordCode endpoint.Endpoint
		ResetPassword             endpoint.Endpoint
		VerifyAccount             endpoint.Endpoint
	}

	authService interface {
		Login(ctx context.Context, email, password string) (string, error)
		Logout(ctx context.Context) error
		SignUp(ctx context.Context, email, password, username string) error
		ForgotPassword(ctx context.Context, email string) error
		ValidateResetPasswordCode(ctx context.Context, email, otp string) (uuid.UUID, error)
		ResetPassword(ctx context.Context, email, password, otp string) error
		VerifyAccount(ctx context.Context, userID uuid.UUID, otp string) error
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
		Username string `json:"user_name" validate:"required"`
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

	// VerifyAccountRequest struct
	VerifyAccountRequest struct {
		UserID string `json:"user_id,omitempty" validate:"required"`
		OTP    string `json:"otp" validate:"required"`
	}

	// ValidateResetPasswordCodeRequest struct
	ValidateResetPasswordCodeRequest struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp" validate:"required"`
	}
)

// MakeEndpoints ...
func MakeEndpoints(as authService, jwtMdw endpoint.Middleware, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		Login:                     MakeLoginEndpoint(as, validateFunc),
		Logout:                    jwtMdw(MakeLogoutEndpoint(as, validateFunc)),
		SignUp:                    MakeSignUpEndpoint(as, validateFunc),
		ForgotPassword:            MakeForgotPasswordEndpoint(as, validateFunc),
		ValidateResetPasswordCode: MakeValidateResetPasswordCodeEndpoint(as, validateFunc),
		ResetPassword:             MakeResetPasswordEndpoint(as, validateFunc),
		VerifyAccount:             MakeVerifyAccountEndpoint(as, validateFunc),
	}

	if len(m) > 0 {
		for _, mdw := range m {
			e.Login = mdw(e.Login)
			e.Logout = mdw(e.Logout)
			e.SignUp = mdw(e.SignUp)
			e.ForgotPassword = mdw(e.ForgotPassword)
			e.ValidateResetPasswordCode = mdw(e.ValidateResetPasswordCode)
			e.ResetPassword = mdw(e.ResetPassword)
			e.VerifyAccount = mdw(e.VerifyAccount)
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
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		if err := s.Logout(ctx); err != nil {
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

		if err := s.SignUp(ctx, req.Email, req.Password, req.Username); err != nil {
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

		if err := s.ForgotPassword(ctx, req.Email); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// MakeValidateResetPasswordCodeEndpoint ...
func MakeValidateResetPasswordCodeEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateResetPasswordCodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		if _, err := s.ValidateResetPasswordCode(ctx, req.Email, req.OTP); err != nil {
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

		if err := s.ResetPassword(ctx, req.Email, req.Password, req.OTP); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// MakeVerifyAccountEndpoint ...
func MakeVerifyAccountEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyAccountRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return nil, err
		}

		if err := s.VerifyAccount(ctx, userID, req.OTP); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

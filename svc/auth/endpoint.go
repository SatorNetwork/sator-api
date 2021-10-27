package auth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/validator"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints struct
	Endpoints struct {
		Login        endpoint.Endpoint
		Logout       endpoint.Endpoint
		SignUp       endpoint.Endpoint
		RefreshToken endpoint.Endpoint

		ForgotPassword            endpoint.Endpoint
		ValidateResetPasswordCode endpoint.Endpoint
		ResetPassword             endpoint.Endpoint
		ChangePassword            endpoint.Endpoint

		VerifyAccount endpoint.Endpoint

		RequestChangeEmail endpoint.Endpoint
		UpdateEmail        endpoint.Endpoint
		UpdateUsername     endpoint.Endpoint

		RequestDestroyAccount endpoint.Endpoint
		VerifyDestroyCode     endpoint.Endpoint
		DestroyAccount        endpoint.Endpoint
		IsVerified            endpoint.Endpoint
		ResendOTP             endpoint.Endpoint
	}

	authService interface {
		Login(ctx context.Context, email, password string) (string, error)
		Logout(ctx context.Context, tid string) error
		SignUp(ctx context.Context, email, password, username string) (string, error)
		RefreshToken(ctx context.Context, uid uuid.UUID, username, tid string) (string, error)

		ForgotPassword(ctx context.Context, email string) error
		ValidateResetPasswordCode(ctx context.Context, email, otp string) (uuid.UUID, error)
		ResetPassword(ctx context.Context, email, password, otp string) error
		ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error

		VerifyAccount(ctx context.Context, userID uuid.UUID, otp string) error

		RequestChangeEmail(ctx context.Context, userID uuid.UUID, email string) error
		ValidateChangeEmailCode(ctx context.Context, userID uuid.UUID, email, otp string) error
		UpdateEmail(ctx context.Context, userID uuid.UUID, email, otp string) error
		UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error

		RequestDestroyAccount(ctx context.Context, uid uuid.UUID) error
		ValidateDestroyAccountCode(ctx context.Context, uid uuid.UUID, otp string) error
		DestroyAccount(ctx context.Context, uid uuid.UUID, otp string) error
		IsVerified(ctx context.Context, userID uuid.UUID) (bool, error)
		ResendOTP(ctx context.Context, userID uuid.UUID) error
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

	// RefreshTokenRequest struct
	RefreshTokenRequest struct {
		UserID  string `json:"user_id,omitempty" validate:"required"`
		TokenID string `json:"token_id,omitempty" validate:"required"`
	}

	// SignUpRequest struct
	SignUpRequest struct {
		Email    string `json:"email" validate:"required,email,lte=150"`
		Password string `json:"password" validate:"required,gte=8,lte=50"`
		Username string `json:"username" validate:"required,alphanum,gte=5,lte=50"`
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

	// ChangePasswordRequest struct
	ChangePasswordRequest struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	}

	// VerifyOTPRequest struct
	VerifyOTPRequest struct {
		OTP string `json:"otp" validate:"required"`
	}

	// ValidateResetPasswordCodeRequest struct
	ValidateResetPasswordCodeRequest struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp" validate:"required"`
	}

	RequestChangeEmailRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	ValidateChangeEmailCodeRequest struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp" validate:"required"`
	}

	UpdateEmailRequest struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp" validate:"required"`
	}

	UpdateUsernameRequest struct {
		Username string `json:"username" validate:"required"`
	}
)

// MakeEndpoints ...
func MakeEndpoints(as authService, jwtMdw endpoint.Middleware, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		Login:        MakeLoginEndpoint(as, validateFunc),
		Logout:       jwtMdw(MakeLogoutEndpoint(as)),
		SignUp:       MakeSignUpEndpoint(as, validateFunc),
		RefreshToken: jwtMdw(MakeRefreshTokenEndpoint(as, validateFunc)),

		ForgotPassword:            MakeForgotPasswordEndpoint(as, validateFunc),
		ValidateResetPasswordCode: MakeValidateResetPasswordCodeEndpoint(as, validateFunc),
		ResetPassword:             MakeResetPasswordEndpoint(as, validateFunc),
		ChangePassword:            MakeChangePasswordEndpoint(as, validateFunc),
		VerifyAccount:             jwtMdw(MakeVerifyAccountEndpoint(as, validateFunc)),
		IsVerified:                jwtMdw(MakeIsVerifiedEndpoint(as)),
		ResendOTP:                 jwtMdw(MakeResendOTPEndpoint(as)),

		RequestChangeEmail: jwtMdw(MakeRequestChangeEmailEndpoint(as, validateFunc)),
		UpdateEmail:        jwtMdw(MakeUpdateEmailEndpoint(as, validateFunc)),
		UpdateUsername:     jwtMdw(MakeUpdateUsernameEndpoint(as, validateFunc)),

		RequestDestroyAccount: jwtMdw(MakeRequestDestroyAccount(as, validateFunc)),
		VerifyDestroyCode:     jwtMdw(MakeVerifyDestroyEndpoint(as, validateFunc)),
		DestroyAccount:        jwtMdw(MakeDestroyAccountEndpoint(as, validateFunc)),
	}

	if len(m) > 0 {
		for _, mdw := range m {
			e.Login = mdw(e.Login)
			e.Logout = mdw(e.Logout)
			e.SignUp = mdw(e.SignUp)
			e.RefreshToken = mdw(e.RefreshToken)

			e.ForgotPassword = mdw(e.ForgotPassword)
			e.ValidateResetPasswordCode = mdw(e.ValidateResetPasswordCode)
			e.ResetPassword = mdw(e.ResetPassword)
			e.ChangePassword = mdw(e.ChangePassword)

			e.VerifyAccount = mdw(e.VerifyAccount)
			e.IsVerified = mdw(e.IsVerified)
			e.ResendOTP = mdw(e.ResendOTP)

			e.RequestChangeEmail = mdw(e.RequestChangeEmail)
			e.UpdateEmail = mdw(e.UpdateEmail)
			e.UpdateUsername = mdw(e.UpdateUsername)

			e.RequestDestroyAccount = mdw(e.RequestDestroyAccount)
			e.VerifyDestroyCode = mdw(e.VerifyDestroyCode)
			e.DestroyAccount = mdw(e.DestroyAccount)
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

		token, err := s.Login(
			ctx,
			strings.ToLower(strings.TrimSpace(req.Email)),
			strings.TrimSpace(req.Password),
		)
		if err != nil {
			if errors.Is(err, ErrInvalidCredentials) {
				return nil, validator.NewValidationError(url.Values{
					"email":    []string{"Invalid email address"},
					"password": []string{"Invalid password"},
				})
			}
			return nil, err
		}

		return AccessToken{Token: token}, nil
	}
}

// MakeLogoutEndpoint ...
func MakeLogoutEndpoint(s authService) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		tid, err := jwt.TokenIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}
		if err := s.Logout(ctx, tid.String()); err != nil {
			return nil, err
		}
		return true, nil
	}
}

// MakeSignUpEndpoint ...
func MakeSignUpEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		token, err := s.SignUp(
			ctx,
			strings.ToLower(strings.TrimSpace(req.Email)),
			strings.TrimSpace(req.Password),
			strings.TrimSpace(req.Username),
		)
		if err != nil {
			return nil, err
		}

		return AccessToken{Token: token}, nil
	}
}

// MakeRefreshTokenEndpoint ...
func MakeRefreshTokenEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		username, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get username: %w", err)
		}

		tid, err := jwt.TokenIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		token, err := s.RefreshToken(ctx, uid, username, tid.String())
		if err != nil {
			return nil, err
		}

		return AccessToken{Token: token}, nil
	}
}

// MakeForgotPasswordEndpoint ...
func MakeForgotPasswordEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ForgotPasswordRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		// normalize email address
		req.Email = strings.ToLower(req.Email)

		if err := s.ForgotPassword(ctx, req.Email); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeValidateResetPasswordCodeEndpoint ...
func MakeValidateResetPasswordCodeEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateResetPasswordCodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		// normalize email address
		req.Email = strings.ToLower(req.Email)

		if _, err := s.ValidateResetPasswordCode(ctx, req.Email, req.OTP); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeResetPasswordEndpoint ...
func MakeResetPasswordEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ResetPasswordRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		// normalize email address
		req.Email = strings.ToLower(req.Email)

		if err := s.ResetPassword(ctx, req.Email, req.Password, req.OTP); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeChangePasswordEndpoint ...
func MakeChangePasswordEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		req := request.(ChangePasswordRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		if err := s.ChangePassword(ctx, uid, req.OldPassword, req.NewPassword); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeVerifyAccountEndpoint ...
func MakeVerifyAccountEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyOTPRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.VerifyAccount(ctx, uid, req.OTP); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeRequestChangeEmailEndpoint ...
func MakeRequestChangeEmailEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RequestChangeEmailRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		// normalize email address
		req.Email = strings.ToLower(req.Email)

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.RequestChangeEmail(ctx, uid, req.Email); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeValidateChangeEmailCodeEndpoint ...
func MakeValidateChangeEmailCodeEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateChangeEmailCodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		// normalize email address
		req.Email = strings.ToLower(req.Email)

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.ValidateChangeEmailCode(ctx, uid, req.Email, req.OTP); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeUpdateEmailEndpoint ...
func MakeUpdateEmailEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateEmailRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		// normalize email address
		req.Email = strings.ToLower(req.Email)

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.UpdateEmail(ctx, uid, req.Email, req.OTP); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeUpdateUsernameEndpoint ...
func MakeUpdateUsernameEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUsernameRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		userName, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get username: %w", err)
		}

		if strings.ToLower(userName) == strings.ToLower(req.Username) {
			return false, nil
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.UpdateUsername(ctx, uid, req.Username); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeRequestDestroyAccount ...
func MakeRequestDestroyAccount(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.RequestDestroyAccount(ctx, uid); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeVerifyDestroyEndpoint ...
func MakeVerifyDestroyEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyOTPRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.ValidateDestroyAccountCode(ctx, uid, req.OTP); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeDestroyAccountEndpoint ...
func MakeDestroyAccountEndpoint(s authService, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyOTPRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.ValidateDestroyAccountCode(ctx, uid, req.OTP); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeIsVerifiedEndpoint ...
func MakeIsVerifiedEndpoint(s authService) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		isVerified, err := s.IsVerified(ctx, uid)
		if err != nil {
			return nil, err
		}

		return isVerified, nil
	}
}

// MakeResendOTPEndpoint ...
func MakeResendOTPEndpoint(s authService) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.ResendOTP(ctx, uid); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

package auth

import (
	"context"

	"github.com/SatorNetwork/sator-api/internal/validator"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints struct
	// TODO: add missed endpoints
	Endpoints struct {
		Login endpoint.Endpoint
	}

	authService interface {
		Login(ctx context.Context, email, password string) (string, error)
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
)

// MakeEndpoints ...
// TODO: add missed endpoints
func MakeEndpoints(as authService, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		Login: MakeLoginEndpoint(as, validateFunc),
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

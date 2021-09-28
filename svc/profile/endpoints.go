package profile

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetProfile   endpoint.Endpoint
		UpdateAvatar endpoint.Endpoint
	}

	service interface {
		GetProfileByUserID(ctx context.Context, userID uuid.UUID, username string) (interface{}, error)
		UpdateAvatar(ctx context.Context, uid uuid.UUID, avatar string) error
	}

	// UpdateAvatarRequest struct
	UpdateAvatarRequest struct {
		Avatar string `json:"avatar,omitempty" validate:"required,gt=0"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetProfile:   MakeGetProfileEndpoint(s),
		UpdateAvatar: MakeUpdateAvatarEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetProfile = mdw(e.GetProfile)
			e.UpdateAvatar = mdw(e.UpdateAvatar)
		}
	}

	return e
}

// MakeGetProfileEndpoint ...
func MakeGetProfileEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		username, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get username: %w", err)
		}

		resp, err := s.GetProfileByUserID(ctx, uid, username)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeUpdateAvatarEndpoint ...
func MakeUpdateAvatarEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(UpdateAvatarRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		err = s.UpdateAvatar(ctx, uid, req.Avatar)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

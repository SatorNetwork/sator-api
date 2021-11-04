package quiz

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/rbac"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetQuizLink endpoint.Endpoint
	}

	service interface {
		GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (interface{}, error)
	}

	// ConnectionURL struct
	ConnectionURL struct {
		PlayURL string `json:"play_url"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetQuizLink: MakeGetQuizLinkEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetQuizLink = mdw(e.GetQuizLink)
		}
	}

	return e
}

// MakeGetQuizLinkEndpoint ...
func MakeGetQuizLinkEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		username, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get username: %w", err)
		}

		challengeID, err := uuid.Parse(req.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		playURL, err := s.GetQuizLink(ctx, uid, username, challengeID)
		if err != nil {
			return nil, err
		}

		return ConnectionURL{
			PlayURL: fmt.Sprintf("%v", playURL),
		}, nil
	}
}

package quiz_v2

import (
	"context"
	"fmt"
	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/rbac"
	challenge_service "github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetQuizLink      endpoint.Endpoint
		GetChallengeById endpoint.Endpoint
	}

	service interface {
		GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (*GetQuizLinkResponse, error)
		GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID) (challenge_service.Challenge, error)
	}

	GetQuizLinkResponse struct {
		BaseQuizWSURL   string `json:"base_quiz_ws_url"`
		BaseQuizURL     string `json:"base_quiz_url"`
		RecvMessageSubj string `json:"recv_message_subj"`
		SendMessageSubj string `json:"send_message_subj"`
		UserID          string `json:"user_id"`
		ServerPublicKey string `json:"server_public_key"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetQuizLink:      MakeGetQuizLinkEndpoint(s),
		GetChallengeById: MakeGetChallengeByIdEndpoint(s),
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
		//		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
		//			return nil, err
		//		}

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

		resp, err := s.GetQuizLink(ctx, uid, username, challengeID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetChallengeByIdEndpoint ...
func MakeGetChallengeByIdEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		userID, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "can't get userid from context")
		}

		challengeID, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, errors.Wrap(err, "can't parse challenge id")
		}

		resp, err := s.GetChallengeByID(ctx, challengeID, userID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

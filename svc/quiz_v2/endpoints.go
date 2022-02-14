package quiz_v2

import (
	"context"
	"fmt"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/SatorNetwork/sator-api/internal/utils"

	"github.com/SatorNetwork/sator-api/svc/challenge"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetQuizLink       endpoint.Endpoint
		GetChallenges     endpoint.Endpoint
		GetFillingQuizzes endpoint.Endpoint
	}

	service interface {
		GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (*GetQuizLinkResponse, error)
		GetChallenges(ctx context.Context, limit, offset int32) (*GetChallengesResponse, error)
		GetFillingQuizzes(ctx context.Context) (*GetFillingQuizzes, error)
	}

	GetQuizLinkResponse struct {
		BaseQuizWSURL   string `json:"base_quiz_ws_url"`
		BaseQuizURL     string `json:"base_quiz_url"`
		RecvMessageSubj string `json:"recv_message_subj"`
		SendMessageSubj string `json:"send_message_subj"`
		UserID          string `json:"user_id"`
		ServerPublicKey string `json:"server_public_key"`
	}

	GetChallengesRequest struct {
		utils.PaginationRequest
	}

	GetChallengesResponse struct {
		Challenges []challenge.Challenge
	}

	GetFillingQuizzes struct {
		playersInRooms map[int32]room.Room
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetQuizLink:       MakeGetQuizLinkEndpoint(s),
		GetChallenges:     MakeGetChallengesEndpoint(s, validateFunc),
		GetFillingQuizzes: MakeGetFillingQuizzesEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetQuizLink = mdw(e.GetQuizLink)
			e.GetChallenges = mdw(e.GetChallenges)
			e.GetFillingQuizzes = mdw(e.GetFillingQuizzes)
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

// MakeGetChallengesEndpoint ...
func MakeGetChallengesEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
		//			return nil, err
		//		}

		req := request.(GetChallengesRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetChallenges(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetFillingQuizzesEndpoint ...
func MakeGetFillingQuizzesEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		//		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
		//			return nil, err
		//		}

		resp, err := s.GetFillingQuizzes(ctx)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

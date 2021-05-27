package questions

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of question service
	Endpoints struct {
		GetQuestion endpoint.Endpoint
	}

	service interface {
		GetQuestionByID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetQuestionByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetQuestion: MakeGetQuestionEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetQuestion = mdw(e.GetQuestion)
		}
	}

	return e
}

// MakeGetQuestionEndpoint ...
func MakeGetQuestionEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user question id: %w", err)
		}

		resp, err := s.GetQuestionByID(ctx, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

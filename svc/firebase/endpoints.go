package firebase

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/validator"
)

type (
	Endpoints struct {
		RegisterToken endpoint.Endpoint
	}

	service interface {
		RegisterToken(ctx context.Context, userID uuid.UUID, req *RegisterTokenRequest) (*Empty, error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		RegisterToken: MakeRegisterTokenEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoint
	if len(m) > 0 {
		for _, mdw := range m {
			e.RegisterToken = mdw(e.RegisterToken)
		}
	}

	return e
}

func MakeRegisterTokenEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		userID, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		typedReq, ok := req.(*RegisterTokenRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to register-token-request")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		_, err = s.RegisterToken(ctx, userID, typedReq)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

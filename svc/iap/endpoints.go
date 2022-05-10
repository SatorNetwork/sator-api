package iap

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
		RegisterInAppPurchase endpoint.Endpoint
	}

	service interface {
		RegisterInAppPurchase(ctx context.Context, userID uuid.UUID, req *RegisterInAppPurchaseRequest) (*Empty, error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		RegisterInAppPurchase: MakeRegisterInAppPurchaseEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoint
	if len(m) > 0 {
		for _, mdw := range m {
			e.RegisterInAppPurchase = mdw(e.RegisterInAppPurchase)
		}
	}

	return e
}

func MakeRegisterInAppPurchaseEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		userID, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		typedReq, ok := req.(*RegisterInAppPurchaseRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to register-in-app-purchase-request")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		_, err = s.RegisterInAppPurchase(ctx, userID, typedReq)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

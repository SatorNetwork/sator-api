package flags

import (
	"context"
	"fmt"
	"github.com/SatorNetwork/sator-api/lib/rbac"

	"github.com/go-kit/kit/endpoint"

	"github.com/SatorNetwork/sator-api/svc/flags/repository"
)

type (
	// Endpoints collection of NFT service
	Endpoints struct {
		UpdateFlag endpoint.Endpoint
		GetFlags   endpoint.Endpoint
	}

	service interface {
		UpdateFlag(ctx context.Context, flag *repository.Flag) (*repository.Flag, error)
		GetFlags(ctx context.Context) ([]repository.Flag, error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		UpdateFlag: MakeUpdateFlagEndpoint(s),
		GetFlags:   MakeGetFlagsEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetFlags = mdw(e.GetFlags)
			e.UpdateFlag = mdw(e.UpdateFlag)
		}
	}

	return e
}

func MakeUpdateFlagEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		req, ok := request.(repository.Flag)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: Flag, got: %T", request)
		}

		return s.UpdateFlag(ctx, &req)
	}
}

func MakeGetFlagsEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		return s.GetFlags(ctx)
	}
}

package wallet

import (
	"context"
	"fmt"
	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetBalance          endpoint.Endpoint
		GetListTransactionsByWalletID endpoint.Endpoint
	}

	service interface {
		GetBalance(ctx context.Context, uid uuid.UUID) (interface{}, error)
		GetListTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) (_ interface{}, err error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetBalance: MakeGetBalanceEndpoint(s),
		GetListTransactionsByWalletID: MakeGetListTransactionsByWalletIDEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetBalance = mdw(e.GetBalance)
			e.GetListTransactionsByWalletID = mdw(e.GetListTransactionsByWalletID)
		}
	}

	return e
}

// MakeGetBalanceEndpoint ...
func MakeGetBalanceEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.GetBalance(ctx, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeGetListTransactionsByWalletIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		walletUID, err := uuid.Parse(req.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get wallet id: %w", err)
		}

		resp, err := s.GetListTransactionsByWalletID(ctx, walletUID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
package rewards

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		ClaimRewards     endpoint.Endpoint
		GetRewardsWallet endpoint.Endpoint
		GetTransactions  endpoint.Endpoint
	}

	service interface {
		ClaimRewards(ctx context.Context, uid uuid.UUID) (ClaimRewardsResult, error)
		GetRewardsWallet(ctx context.Context, userID, walletID uuid.UUID) (wallet.Wallet, error)
		GetTransactions(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (wallet.Transactions, error)
	}

	// GetTransactionsRequest struct
	GetTransactionsRequest struct {
		WalletID string
		PaginationRequest
	}

	// PaginationRequest struct
	PaginationRequest struct {
		Page         int32 `json:"page,omitempty" validate:"number,gte=0"`
		ItemsPerPage int32 `json:"items_per_page,omitempty" validate:"number,gte=0"`
	}
)

// Limit of items
func (r PaginationRequest) Limit() int32 {
	if r.ItemsPerPage > 0 {
		return r.ItemsPerPage
	}
	return 20
}

// Offset items
func (r PaginationRequest) Offset() int32 {
	if r.Page > 1 {
		return (r.Page - 1) * r.Limit()
	}
	return 0
}

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		ClaimRewards:     MakeClaimRewardsEndpoint(s),
		GetRewardsWallet: MakeGetRewardsWalletEndpoint(s),
		GetTransactions:  MakeGetTransactionsEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.ClaimRewards = mdw(e.ClaimRewards)
			e.GetRewardsWallet = mdw(e.GetRewardsWallet)
			e.GetTransactions = mdw(e.GetTransactions)
		}
	}

	return e
}

// MakeClaimRewardsEndpoint ...
func MakeClaimRewardsEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		rewards, err := s.ClaimRewards(ctx, uid)
		if err != nil {
			return nil, err
		}

		return rewards, nil
	}
}

// MakeGetRewardsWalletEndpoint ...
func MakeGetRewardsWalletEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		walletID, err := uuid.Parse(req.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get wallet id: %w", err)
		}

		wallet, err := s.GetRewardsWallet(ctx, uid, walletID)
		if err != nil {
			return nil, err
		}

		return wallet, nil
	}
}

// MakeGetTransactionsEndpoint ...
func MakeGetTransactionsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(GetTransactionsRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		walletID, err := uuid.Parse(req.WalletID)
		if err != nil {
			return nil, fmt.Errorf("could not parse wallet id: %w", err)
		}

		transactions, err := s.GetTransactions(ctx, uid, walletID, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return transactions, nil
	}
}

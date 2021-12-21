package wallet

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/rbac"
	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		CreateTransfer                endpoint.Endpoint
		ConfirmTransfer               endpoint.Endpoint
		GetWallets                    endpoint.Endpoint
		GetWalletByID                 endpoint.Endpoint
		GetListTransactionsByWalletID endpoint.Endpoint
		GetStake                      endpoint.Endpoint
		SetStake                      endpoint.Endpoint
	}

	service interface {
		GetListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (_ Transactions, err error)
		GetWallets(ctx context.Context, uid uuid.UUID) (Wallets, error)
		GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (Wallet, error)
		CreateTransfer(ctx context.Context, senderWalletID uuid.UUID, recipientAddr, asset string, amount float64) (PreparedTransferTransaction, error)
		ConfirmTransfer(ctx context.Context, senderWalletID uuid.UUID, tx string) error
		GetStake(ctx context.Context, walletID uuid.UUID) (Stake, error)
		SetStake(ctx context.Context, walletID uuid.UUID, amount float64) (bool, error)
	}

	CreateTransferRequest struct {
		SenderWalletID   string  `json:"-"`
		RecipientAddress string  `json:"recipient_address" validate:"required"`
		Amount           float64 `json:"amount" validate:"required,number,gt=0"`
		Asset            string  `json:"asset,omitempty"`
	}

	ConfirmTransferRequest struct {
		SenderWalletID  string `json:"-"`
		TransactionHash string `json:"tx_hash"`
	}

	// GetListTransactionsByWalletIDRequest struct
	GetListTransactionsByWalletIDRequest struct {
		WalletID string `json:"wallet_id" validate:"required,uuid"`
		utils.PaginationRequest
	}

	// SetStakeRequest struct
	SetStakeRequest struct {
		Amount   float64 `json:"amount" validate:"required,number,gt=0"`
		WalletID string  `json:"wallet_id" validate:"required,uuid"`
	}
)

func MakeEndpoints(s service, kycMdw endpoint.Middleware, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetWallets:                    MakeGetWalletsEndpoint(s),
		GetWalletByID:                 MakeGetWalletByIDEndpoint(s),
		GetListTransactionsByWalletID: MakeGetListTransactionsByWalletIDEndpoint(s, validateFunc),
		CreateTransfer:                kycMdw(MakeCreateTransferRequestEndpoint(s, validateFunc)),
		ConfirmTransfer:               kycMdw(MakeConfirmTransferRequestEndpoint(s, validateFunc)),
		SetStake:                      MakeSetStakeEndpoint(s, validateFunc),
		GetStake:                      MakeGetStakeEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetWallets = mdw(e.GetWallets)
			e.GetWalletByID = mdw(e.GetWalletByID)
			e.GetListTransactionsByWalletID = mdw(e.GetListTransactionsByWalletID)
			e.CreateTransfer = mdw(e.CreateTransfer)
			e.ConfirmTransfer = mdw(e.ConfirmTransfer)
			e.SetStake = mdw(e.SetStake)
			e.GetStake = mdw(e.GetStake)
		}
	}

	return e
}

func MakeGetListTransactionsByWalletIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(GetListTransactionsByWalletIDRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		walletUID, err := uuid.Parse(req.WalletID)
		if err != nil {
			return nil, fmt.Errorf("could not get wallet id: %w", err)
		}

		transactions, err := s.GetListTransactionsByWalletID(ctx, uid, walletUID, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return transactions, nil
	}
}

func MakeGetWalletsEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		wallets, err := s.GetWallets(ctx, uid)
		if err != nil {
			return nil, err
		}

		return wallets, nil
	}
}

func MakeGetWalletByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		walletID, err := uuid.Parse(req.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get wallet id: %w", err)
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		w, err := s.GetWalletByID(ctx, uid, walletID)
		if err != nil {
			return nil, err
		}

		return w, nil
	}
}

func MakeCreateTransferRequestEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(CreateTransferRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		walletID, err := uuid.Parse(req.SenderWalletID)
		if err != nil {
			return nil, fmt.Errorf("invalid sender wallet id: %w", err)
		}

		txInfo, err := s.CreateTransfer(ctx, walletID, req.RecipientAddress, req.Asset, req.Amount)
		if err != nil {
			return nil, err
		}

		return txInfo, nil
	}
}

func MakeConfirmTransferRequestEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
		// 	return nil, err
		// }

		req := request.(ConfirmTransferRequest)
		if err := v(req); err != nil {
			return false, err
		}

		walletID, err := uuid.Parse(req.SenderWalletID)
		if err != nil {
			return nil, fmt.Errorf("invalid sender wallet id: %w", err)
		}

		if err := s.ConfirmTransfer(ctx, walletID, req.TransactionHash); err != nil {
			return false, err
		}

		return true, nil
	}
}

func MakeSetStakeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(SetStakeRequest)
		if err := v(req); err != nil {
			return false, err
		}

		walletID, err := uuid.Parse(req.WalletID)
		if err != nil {
			return nil, fmt.Errorf("invalid wallet id: %w", err)
		}

		result, err := s.SetStake(ctx, walletID, req.Amount)
		if err != nil {
			return false, err
		}

		return result, nil
	}
}

func MakeGetStakeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		walletID, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get wallet id: %w", err)
		}

		stake, err := s.GetStake(ctx, walletID)
		if err != nil {
			return nil, err
		}

		return stake, nil
	}
}

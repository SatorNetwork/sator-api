package wallet

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/solana/client"
	"github.com/SatorNetwork/sator-api/lib/utils"
	"github.com/SatorNetwork/sator-api/lib/validator"
	"github.com/dmitrymomot/go-env"
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
		Unstake                       endpoint.Endpoint
		PossibleMultiplier            endpoint.Endpoint
		GetStakeLevels                endpoint.Endpoint
	}

	service interface {
		GetListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (_ Transactions, err error)
		GetWallets(ctx context.Context, uid uuid.UUID) (Wallets, error)
		GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (Wallet, error)
		CreateTransfer(ctx context.Context, senderWalletID uuid.UUID, recipientAddr, asset string, amount float64) (PreparedTransferTransaction, error)
		ConfirmTransfer(ctx context.Context, senderWalletID uuid.UUID, tx string) error
		GetStake(ctx context.Context, userID uuid.UUID) (Stake, error)
		SetStake(ctx context.Context, userID, walletID uuid.UUID, duration int64, amount float64) (bool, error)
		Unstake(ctx context.Context, userID, walletID uuid.UUID) error
		PossibleMultiplier(ctx context.Context, additionalAmount float64, userID, walletID uuid.UUID) (int32, error)
		GetEnabledStakeLevelsList(ctx context.Context, userID uuid.UUID) ([]StakeLevel, error)
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
		Duration int64   `json:"duration"`
	}

	// UnstakeRequest struct
	UnstakeRequest struct {
		WalletID string `json:"wallet_id" validate:"required,uuid"`
	}

	// PossibleMultiplierRequest struct
	PossibleMultiplierRequest struct {
		Amount   float64 `json:"amount" validate:"required,number,gt=0"`
		WalletID string  `json:"wallet_id" validate:"required,uuid"`
	}
)

var StakeDuration = env.GetInt("SMART_CONTRACT_STAKE_DURATION", 0)

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
		Unstake:                       MakeUnstakeEndpoint(s, validateFunc),
		PossibleMultiplier:            MakePossibleMultiplierEndpoint(s, validateFunc),
		GetStakeLevels:                MakeGetStakeLevelsEndpoint(s),
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
			e.Unstake = mdw(e.Unstake)
			e.PossibleMultiplier = mdw(e.PossibleMultiplier)
			e.GetStakeLevels = mdw(e.GetStakeLevels)
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

		if err := validateSolanaWalletAddr("recipient_address", req.RecipientAddress); err != nil {
			return nil, err
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

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(SetStakeRequest)
		if err := v(req); err != nil {
			return false, err
		}

		walletID, err := uuid.Parse(req.WalletID)
		if err != nil {
			return nil, fmt.Errorf("invalid wallet id: %w", err)
		}

		if req.Duration == 0 {
			req.Duration = int64(StakeDuration)
		}

		result, err := s.SetStake(ctx, uid, walletID, req.Duration, req.Amount)
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

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		stake, err := s.GetStake(ctx, uid)
		if err != nil {
			return nil, err
		}

		return stake, nil
	}
}

func MakeUnstakeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return false, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return false, fmt.Errorf("could not get user profile id: %w", err)
		}

		walletID, err := uuid.Parse(request.(string))
		if err != nil {
			return false, fmt.Errorf("could not get wallet id: %w", err)
		}

		err = s.Unstake(ctx, uid, walletID)
		if err != nil {
			return false, err
		}

		return true, nil
	}
}

func MakePossibleMultiplierEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(PossibleMultiplierRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		walletID, err := uuid.Parse(req.WalletID)
		if err != nil {
			return nil, fmt.Errorf("invalid wallet id: %w", err)
		}

		multiplier, err := s.PossibleMultiplier(ctx, req.Amount, uid, walletID)
		if err != nil {
			return nil, err
		}

		return multiplier, nil
	}
}

func MakeGetStakeLevelsEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		levels, err := s.GetEnabledStakeLevelsList(ctx, uid)
		if err != nil {
			return nil, err
		}

		return levels, nil
	}
}

func validateSolanaWalletAddr(fieldName, addr string) error {
	if err := client.ValidateSolanaWalletAddr(addr); err != nil {
		log.Printf("invalid solana wallet address=%s, error: %v", addr, err)
		return validator.NewValidationError(url.Values{
			fieldName: []string{"invalid solana wallet address"},
		})
	}

	return nil
}

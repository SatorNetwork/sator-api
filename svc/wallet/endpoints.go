package wallet

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		Transfer        endpoint.Endpoint
		CreateTransfer  endpoint.Endpoint
		ConfirmTransfer endpoint.Endpoint

		GetWallets                    endpoint.Endpoint
		GetWalletByID                 endpoint.Endpoint
		GetListTransactionsByWalletID endpoint.Endpoint
	}

	service interface {
		GetListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (_ Transactions, err error)
		GetWallets(ctx context.Context, uid uuid.UUID) (Wallets, error)
		GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (Wallet, error)

		Transfer(ctx context.Context, senderPrivateKey, recipientPK string, amount float64) (tx string, err error)
		CreateTransfer(ctx context.Context, senderWalletID uuid.UUID, recipientAddr string, amount float64) (PreparedTransaction, error)
		ConfirmTransfer(ctx context.Context, senderWalletID uuid.UUID, tx string) error
	}

	CreateTransferRequest struct {
		SenderWalletID   string  `json:"-"`
		RecipientAddress string  `json:"recipient_address" validate:"required"`
		Amount           float64 `json:"amount" validate:"required,number,gt=0"`
	}

	ConfirmTransferRequest struct {
		SenderWalletID  string `json:"-"`
		TransactionHash string `json:"tx_hash"`
	}

	TransferRequest struct {
		SenderPrivateKey string  `json:"sender_private_key" validate:"required"`
		RecipientPK      string  `json:"recipient_pk" validate:"required"`
		Amount           float64 `json:"amount" validate:"required"`
	}

	// GetListTransactionsByWalletIDRequest struct
	GetListTransactionsByWalletIDRequest struct {
		WalletID string `json:"wallet_id" validate:"required,uuid"`
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
		GetWallets:                    MakeGetWalletsEndpoint(s),
		GetWalletByID:                 MakeGetWalletByIDEndpoint(s),
		GetListTransactionsByWalletID: MakeGetListTransactionsByWalletIDEndpoint(s, validateFunc),

		Transfer:        MakeTransferEndpoint(s, validateFunc),
		CreateTransfer:  MakeCreateTransferRequestEndpoint(s, validateFunc),
		ConfirmTransfer: MakeConfirmTransferRequestEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetWallets = mdw(e.GetWallets)
			e.GetWalletByID = mdw(e.GetWalletByID)
			e.GetListTransactionsByWalletID = mdw(e.GetListTransactionsByWalletID)
			e.Transfer = mdw(e.Transfer)
		}
	}

	return e
}

func MakeGetListTransactionsByWalletIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
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
		req := request.(CreateTransferRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		walletID, err := uuid.Parse(req.SenderWalletID)
		if err != nil {
			return nil, fmt.Errorf("invalid sender wallet id: %w", err)
		}

		txInfo, err := s.CreateTransfer(ctx, walletID, req.RecipientAddress, req.Amount)
		if err != nil {
			return nil, err
		}

		return txInfo, nil
	}
}

func MakeConfirmTransferRequestEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
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

func MakeTransferEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TransferRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		tx, err := s.Transfer(ctx, req.SenderPrivateKey, req.RecipientPK, req.Amount)
		if err != nil {
			return nil, err
		}

		return tx, nil
	}
}

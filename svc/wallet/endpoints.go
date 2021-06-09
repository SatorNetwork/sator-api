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
<<<<<<< HEAD
		Transfer                      endpoint.Endpoint
		GetBalance                    endpoint.Endpoint
		GetWallets                    endpoint.Endpoint
=======
		GetBalance          endpoint.Endpoint
>>>>>>> wallets: getListTranscations added
		GetListTransactionsByWalletID endpoint.Endpoint
	}

	service interface {
<<<<<<< HEAD
		GetBalanceWithRewards(ctx context.Context, uid uuid.UUID) (interface{}, error)
		GetListTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) (_ interface{}, err error)
		GetBalanceByUserID(ctx context.Context, userID uuid.UUID) ([]Balance, error)
		Transfer(ctx context.Context, senderPrivateKey, recipientPK string, amount float64) (tx string, err error)
	}

	TransferRequest struct {
		SenderPrivateKey string  `json:"sender_private_key" validate:"required"`
		RecipientPK      string  `json:"recipient_pk" validate:"required"`
		Amount           float64 `json:"amount" validate:"required"`
=======
		GetBalance(ctx context.Context, uid uuid.UUID) (interface{}, error)
		GetListTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) (_ interface{}, err error)
>>>>>>> wallets: getListTranscations added
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
<<<<<<< HEAD
		Transfer:                      MakeTransferEndpoint(s, validateFunc),
		GetBalance:                    MakeGetBalanceEndpoint(s),
		GetWallets:                    MakeGetWalletsEndpoint(s),
=======
		GetBalance: MakeGetBalanceEndpoint(s),
>>>>>>> wallets: getListTranscations added
		GetListTransactionsByWalletID: MakeGetListTransactionsByWalletIDEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.Transfer = mdw(e.Transfer)
			e.GetBalance = mdw(e.GetBalance)
			e.GetListTransactionsByWalletID = mdw(e.GetListTransactionsByWalletID)
<<<<<<< HEAD
			e.GetWallets = mdw(e.GetWallets)
=======
>>>>>>> wallets: getListTranscations added
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

		balance, err := s.GetBalanceWithRewards(ctx, uid)
		if err != nil {
			return nil, err
		}

		return balance, nil
	}
}

func MakeGetListTransactionsByWalletIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		walletUID, err := uuid.Parse(req.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get wallet id: %w", err)
		}

		transactions, err := s.GetListTransactionsByWalletID(ctx, walletUID)
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

		balance, err := s.GetBalanceByUserID(ctx, uid)

		return balance, nil
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
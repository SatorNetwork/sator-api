package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/SatorNetwork/sator-api/lib/utils"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Predefined request query keys
const (
	pageParam         = "page"
	itemsPerPageParam = "items_per_page"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Get("/", httptransport.NewServer(
		e.GetWallets,
		decodeGetWalletsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/sao", httptransport.NewServer(
		e.GetUserWallet,
		decodeGetCurrentUserWalletRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/stake-levels", httptransport.NewServer(
		e.GetStakeLevels,
		decodeGetStakeLevelsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{wallet_id}", httptransport.NewServer(
		e.GetWalletByID,
		decodeGetWalletByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{wallet_id}/transactions", httptransport.NewServer(
		e.GetListTransactionsByWalletID,
		decodeGetListTransactionsByWalletIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{wallet_id}/create-transfer", httptransport.NewServer(
		e.CreateTransfer,
		decodeCreateTransferRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{wallet_id}/confirm-transfer", httptransport.NewServer(
		e.ConfirmTransfer,
		decodeConfirmTransferRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{wallet_id}/stake", httptransport.NewServer(
		e.SetStake,
		decodeSetStakeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{wallet_id}/unstake", httptransport.NewServer(
		e.Unstake,
		decodeUnstakeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{wallet_id}/stake", httptransport.NewServer(
		e.GetStake,
		decodeGetStakeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{wallet_id}/possible-multiplier", httptransport.NewServer(
		e.PossibleMultiplier,
		decodePossibleMultiplierRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetCurrentUserWalletRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetWalletsRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetWalletByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "wallet_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed wallet_id id", ErrInvalidParameter)
	}
	return id, nil
}

func decodeGetListTransactionsByWalletIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetListTransactionsByWalletIDRequest{
		WalletID: chi.URLParam(r, "wallet_id"),
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func decodeCreateTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	senderWalletID := chi.URLParam(r, "wallet_id")
	if senderWalletID == "" {
		return nil, fmt.Errorf("%w: missed wallet_id id", ErrInvalidParameter)
	}
	req.SenderWalletID = senderWalletID
	req.Asset = "SAO" // FIXME: remove hardcode

	return req, nil
}

func decodeConfirmTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ConfirmTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	req.SenderWalletID = chi.URLParam(r, "wallet_id")

	return req, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden, err.Error()
	}

	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, err.Error()
	}

	if errors.Is(err, ErrInvalidParameter) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, ErrTransactionFailed) {
		log.Printf("%v", err)
		return http.StatusInternalServerError, ErrTransactionFailed
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeGetStakeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetStakeLevelsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeSetStakeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SetStakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	req.WalletID = chi.URLParam(r, "wallet_id")

	return req, nil
}

func decodeUnstakeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "wallet_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed wallet_id id", ErrInvalidParameter)
	}
	return id, nil
}

func decodePossibleMultiplierRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req PossibleMultiplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	req.WalletID = chi.URLParam(r, "wallet_id")

	return req, nil
}

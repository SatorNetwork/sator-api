package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
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

	// r.Post("/{wallet_id}/transfer", httptransport.NewServer(
	// 	e.Transfer,
	// 	decodeTransferRequest,
	// 	httpencoder.EncodeResponse,
	// 	options...,
	// ).ServeHTTP)

	return r
}

func decodeGetBalanceRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
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
	id := chi.URLParam(r, "wallet_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed wallet_id id", ErrInvalidParameter)
	}
	return id, nil
}

func decodeTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

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

	return httpencoder.CodeAndMessageFrom(err)
}

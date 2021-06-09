package wallet

import (
	"context"
<<<<<<< HEAD
	"encoding/json"
=======
>>>>>>> wallets: getListTranscations added
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

	r.Get("/balance", httptransport.NewServer(
		e.GetBalance,
		decodeGetBalanceRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

<<<<<<< HEAD
	r.Get("/transactions/{wallet_id}", httptransport.NewServer(
=======
	r.Get("/transactions", httptransport.NewServer(
>>>>>>> wallets: getListTranscations added
		e.GetListTransactionsByWalletID,
		decodeGetListTransactionsByWalletIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

<<<<<<< HEAD
	r.Get("/wallets", httptransport.NewServer(
		e.GetWallets,
		decodeGetWalletsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/transfer", httptransport.NewServer(
		e.Transfer,
		decodeTransferRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

=======
>>>>>>> wallets: getListTranscations added
	return r
}

func decodeGetBalanceRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

<<<<<<< HEAD
func decodeGetWalletsRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

=======
>>>>>>> wallets: getListTranscations added
func decodeGetListTransactionsByWalletIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "wallet_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed qrcode id", ErrInvalidParameter)
	}
	return id, nil
}

<<<<<<< HEAD
func decodeTransferRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

=======
>>>>>>> wallets: getListTranscations added
// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

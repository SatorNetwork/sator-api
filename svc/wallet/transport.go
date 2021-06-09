package wallet

import (
	"context"
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

	r.Get("/transactions", httptransport.NewServer(
		e.GetListTransactionsByWalletID,
		decodeGetListTransactionsByWalletIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetBalanceRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetListTransactionsByWalletIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "wallet_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed qrcode id", ErrInvalidParameter)
	}
	return id, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

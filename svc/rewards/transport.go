package rewards

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
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

	r.Get("/claim", httptransport.NewServer(
		e.ClaimRewards,
		decodeClaimRewardsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/wallet/{wallet_id}", httptransport.NewServer(
		e.GetRewardsWallet,
		decodeGetRewardsWalletRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/wallet/{wallet_id}/transactions", httptransport.NewServer(
		e.GetTransactions,
		decodeGetTransactionsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeClaimRewardsRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetRewardsWalletRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "wallet_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed wallet_id id", ErrInvalidParameter)
	}
	return id, nil
}

func decodeGetTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetTransactionsRequest{
		WalletID: chi.URLParam(r, "wallet_id"),
		PaginationRequest: PaginationRequest{
			Page:         castStrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: castStrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

func castStrToInt32(source string) int32 {
	res, err := strconv.ParseInt(source, 10, 32)
	if err != nil {
		return 0
	}
	return int32(res)
}

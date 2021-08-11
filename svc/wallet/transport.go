package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

	return r
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
		PaginationRequest: PaginationRequest{
			Page:         castStrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: castStrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func decodeCreateTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	req.SenderWalletID = chi.URLParam(r, "wallet_id")
	req.Asset = "SAO" // FIXME: remove hardcode

	log.Printf("\n\nCreateTransferRequest: \n%#v\n\n", req)

	return req, nil
}

func decodeConfirmTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ConfirmTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	req.SenderWalletID = chi.URLParam(r, "wallet_id")
	log.Printf("\n\nConfirmTransferRequest: \n%#v\n\n", req)

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

func castStrToInt32(source string) int32 {
	res, err := strconv.ParseInt(source, 10, 32)
	if err != nil {
		return 0
	}
	return int32(res)
}

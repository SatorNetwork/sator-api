package referrals

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/SatorNetwork/sator-api/internal/utils"

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

	r.Get("/codes", httptransport.NewServer(
		e.GetReferralCodesDataList,
		decodeGetReferralCodesDataListRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/codes/my", httptransport.NewServer(
		e.GetMyReferralCode,
		decodeGetMyReferralCodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/codes", httptransport.NewServer(
		e.AddReferralCodeData,
		decodeAddReferralCodeDataRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/codes/{id}", httptransport.NewServer(
		e.UpdateReferralCodeData,
		decodeUpdateReferralCodeDataRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/codes/{id}", httptransport.NewServer(
		e.DeleteReferralCodeDataByID,
		decodeDeleteReferralCodeDataByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/codes/{id}/referrals", httptransport.NewServer(
		e.GetReferralsWithPaginationByUserID,
		decodeGetReferralsWithPaginationByUserIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/confirm/{code}", httptransport.NewServer(
		e.StoreUserWithValidCode,
		decodeStoreUserWithValidCodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetReferralCodesDataListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return utils.PaginationRequest{
		Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
		ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
	}, nil
}

func decodeGetMyReferralCodeRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeAddReferralCodeDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddReferralCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeUpdateReferralCodeDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateReferralCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed id", ErrInvalidParameter)
	}
	req.ID = id

	return req, nil
}

func decodeDeleteReferralCodeDataByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeStoreUserWithValidCodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	code := chi.URLParam(r, "code")
	if code == "" {
		return nil, fmt.Errorf("%w: missed referral code", ErrInvalidParameter)
	}

	return code, nil
}

func decodeGetReferralsWithPaginationByUserIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return utils.PaginationRequest{
		Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
		ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
	}, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

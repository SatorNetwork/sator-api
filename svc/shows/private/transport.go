package private

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/SatorNetwork/sator-api/internal/rbac"
	"github.com/SatorNetwork/sator-api/internal/utils"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Get("/", httptransport.NewServer(
		e.GetShows,
		decodeGetShowsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) ||
		errors.Is(err, ErrAlreadyReviewed) ||
		errors.Is(err, ErrMaxClaps) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, ErrNotFound) || db.IsNotFoundError(err) {
		return http.StatusNotFound, err.Error()
	}

	if errors.Is(err, rbac.ErrAccessDenied) {
		return http.StatusForbidden, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeGetShowsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetShowsRequest{
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
		},
	}, nil
}

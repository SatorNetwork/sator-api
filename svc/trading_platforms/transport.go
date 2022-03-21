package trading_platforms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/SatorNetwork/sator-api/lib/utils"
)

const (
	pageParam         = "page"
	itemsPerPageParam = "items_per_page"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Get("/", httptransport.NewServer(
		e.GetLinks,
		decodeGetLinksRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/", httptransport.NewServer(
		e.CreateLink,
		decodeCreateLinkRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{id}", httptransport.NewServer(
		e.UpdateLink,
		decodeUpdateLinkRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{id}", httptransport.NewServer(
		e.DeleteLink,
		decodeDeleteLinkRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeCreateLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "could not decode request body")
	}

	return &req, nil
}

func decodeUpdateLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "could not decode request body")
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed link id", ErrInvalidParameter)
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	req.ID = uid

	return &req, nil
}

func decodeDeleteLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed link id", ErrInvalidParameter)
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return &DeleteLinkRequest{
		ID: uid,
	}, nil
}

func decodeGetLinksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return utils.PaginationRequest{
		Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
		ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
	}, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

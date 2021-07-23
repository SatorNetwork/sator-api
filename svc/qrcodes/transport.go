package qrcodes

import (
	"context"
	"errors"
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

	r.Get("/{id}", httptransport.NewServer(
		e.GetDataByQRCodeID,
		decodeGetDataByQRCodeIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/", httptransport.NewServer(
		e.GetQRCodesData,
		decodeGetQRCodesDataRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetDataByQRCodeIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed qrcode id", ErrInvalidParameter)
	}

	return id, nil
}

// codeAndMessageFrom returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrQRCodeExpired) {
		return http.StatusBadRequest, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeGetQRCodesDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return PaginationRequest{
		Page:         castStrToInt32(r.URL.Query().Get(pageParam)),
		ItemsPerPage: castStrToInt32(r.URL.Query().Get(itemsPerPageParam)),
	}, nil
}

func castStrToInt32(source string) int32 {
	res, err := strconv.Atoi(source)
	if err != nil {
		return 0
	}
	return int32(res)
}

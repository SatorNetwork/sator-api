package qrcodes

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

	r.Get("/{id}", httptransport.NewServer(
		e.GetDataByQRCodeID,
		decodeGetDataRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/", httptransport.NewServer(
		e.AddQRCode,
		decodeAddQRCodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{id}", httptransport.NewServer(
		e.DeleteQRCodeByID,
		decodeDeleteQRCodeByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{id}", httptransport.NewServer(
		e.UpdateQRCode,
		decodeUpdateQRCodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

func decodeAddQRCodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddQRCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeDeleteQRCodeByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed qrcode id", ErrInvalidParameter)
	}

	return id, nil

}

func decodeUpdateQRCodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateQRCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	qrcodeID := chi.URLParam(r, "id")
	if qrcodeID == "" {
		return nil, fmt.Errorf("%w: missed qrcode id", ErrInvalidParameter)
	}
	req.ID = qrcodeID

	return req, nil
}

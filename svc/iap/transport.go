package iap

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
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

	r.Post("/register", httptransport.NewServer(
		e.RegisterInAppPurchase,
		decodeRegisterInAppPurchaseRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeRegisterInAppPurchaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RegisterInAppPurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "could not decode request body")
	}

	return &req, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

package flags

import (
	"encoding/json"
	"fmt"
	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/SatorNetwork/sator-api/svc/flags/repository"
	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
	"net/http"
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
		e.GetFlags,
		decodeGetFlagsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/", httptransport.NewServer(
		e.UpdateFlag,
		decodeUpdateFlagRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetFlagsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeUpdateFlagRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req repository.Flag
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

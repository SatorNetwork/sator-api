package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/go-chi/chi/v5"
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

	r.Get("/", httptransport.NewServer(
		e.GetProfile,
		decodeGetProfileRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/avatar", httptransport.NewServer(
		e.UpdateAvatar,
		decodeUpdateAvatarRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetProfileRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeUpdateAvatarRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

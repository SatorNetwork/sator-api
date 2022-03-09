package invitations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"

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

	r.Post("/", httptransport.NewServer(
		e.SendInvitation,
		decodeSendInvitationRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeSendInvitationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SendInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	return req, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) {
		return http.StatusBadRequest, err.Error()
	}
	return httpencoder.CodeAndMessageFrom(err)
}

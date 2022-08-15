package quiz

import (
	"context"
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
func MakeHTTPHandler(e Endpoints, log logger, quizWsHandler http.HandlerFunc) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Get("/{challenge_id}/play", httptransport.NewServer(
		e.GetQuizLink,
		decodeGetQuizLinkRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{challenge_id}/play/{token}", quizWsHandler)

	return r
}

func decodeGetQuizLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "challenge_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed challenge id", ErrInvalidParameter)
	}
	return id, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

package quiz_v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

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

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
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

	r.Get("/challenges/{challenge_id}", httptransport.NewServer(
		e.GetChallengeById,
		decodeGetChallengeByIdRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/challenges/sorted_by_players", httptransport.NewServer(
		e.GetChallengesSortedByPlayers,
		decodeGetChallengesSortedByPlayersRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetQuizLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "challenge_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed challenge id", ErrInvalidParameter)
	}
	return id, nil
}

func decodeGetChallengeByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "challenge_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed challenge id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeGetChallengesSortedByPlayersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return utils.PaginationRequest{
		Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
		ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
	}, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

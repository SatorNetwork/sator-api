package shows

import (
	"context"
	"encoding/json"
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

	r.Get("/", httptransport.NewServer(
		e.GetShows,
		decodeGetShowsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/", httptransport.NewServer(
		e.AddShow,
		decodeAddShowRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{show_id}/challenges", httptransport.NewServer(
		e.GetShowChallenges,
		decodeGetShowChallengesRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{show_id}", httptransport.NewServer(
		e.GetShowByID,
		decodeGetShowByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Patch("/{show_id}", httptransport.NewServer(
		e.UpdateShow,
		decodeUpdateShowRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{show_id}", httptransport.NewServer(
		e.DeleteShowByID,
		decodeDeleteShowByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/filter/{category}", httptransport.NewServer(
		e.GetShowsByCategory,
		decodeGetShowsByCategoryRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/episodes", httptransport.NewServer(
		e.AddEpisode,
		decodeAddEpisodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/episodes/{id}}", httptransport.NewServer(
		e.GetEpisodeByID,
		decodeGetEpisodeByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{show-id}/episodes}", httptransport.NewServer(
		e.GetEpisodesByShowID,
		decodeGetEpisodesByShowIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Patch("/episodes/{id}", httptransport.NewServer(
		e.UpdateEpisode,
		decodeUpdateEpisodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/episodes/{id}", httptransport.NewServer(
		e.DeleteEpisodeByID,
		decodeDeleteEpisodeByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetShowsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return PaginationRequest{
		Page:         castStrToInt32(r.URL.Query().Get(pageParam)),
		ItemsPerPage: castStrToInt32(r.URL.Query().Get(itemsPerPageParam)),
	}, nil
}

func decodeGetShowChallengesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetShowChallengesRequest{
		ShowID: chi.URLParam(r, "show_id"),
		PaginationRequest: PaginationRequest{
			Page:         castStrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: castStrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func decodeGetShowByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "show_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeGetShowsByCategoryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return GetShowsByCategoryRequest{
		Category: chi.URLParam(r, "category"),
		PaginationRequest: PaginationRequest{
			Page:         castStrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: castStrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func castStrToInt32(source string) int32 {
	res, err := strconv.Atoi(source)
	if err != nil {
		return 0
	}

	return int32(res)
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) {
		return http.StatusBadRequest, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeAddShowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddShowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeUpdateShowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateShowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeDeleteShowByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "show_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed show_id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeAddEpisodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddEpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeUpdateEpisodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateEpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	req.ID = chi.URLParam(r, "id")

	return req, nil
}

func decodeDeleteEpisodeByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed episode id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeGetEpisodeByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed episode id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeGetEpisodesByShowIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetEpisodesByShowIDRequest{
		ShowID: chi.URLParam(r, "show_id"),
		PaginationRequest: PaginationRequest{
			Page:         castStrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: castStrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

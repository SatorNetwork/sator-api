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

	r.Get("/home", httptransport.NewServer(
		e.GetShowsHome,
		decodeGetShowsHomeRequest,
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

	r.Put("/{show_id}", httptransport.NewServer(
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

	r.Get("/{show_id}/episodes", httptransport.NewServer(
		e.GetEpisodesByShowID,
		decodeGetEpisodesByShowIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{show_id}/episodes/{episode_id}", httptransport.NewServer(
		e.GetEpisodeByID,
		decodeGetEpisodeByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{show_id}/episodes", httptransport.NewServer(
		e.AddEpisode,
		decodeAddEpisodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{show_id}/episodes/{episode_id}", httptransport.NewServer(
		e.UpdateEpisode,
		decodeUpdateEpisodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{show_id}/episodes/{episode_id}", httptransport.NewServer(
		e.DeleteEpisodeByID,
		decodeDeleteEpisodeByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/categories/{category_id}", httptransport.NewServer(
		e.GetShowCategoryByID,
		decodeGetShowCategoryByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/categories", httptransport.NewServer(
		e.AddShowCategories,
		decodeAddShowCategoriesRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/categories/{category_id}", httptransport.NewServer(
		e.UpdateShowCategory,
		decodeUpdateShowCategoryRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/categories/{category_id}", httptransport.NewServer(
		e.DeleteShowCategoryByID,
		decodeDeleteShowCategoryByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/to-categories/{show_id}", httptransport.NewServer(
		e.DeleteShowToCategoryByShowID,
		decodeDeleteShowToCategoryByShowIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/to-categories/{category_id}/{show_id}", httptransport.NewServer(
		e.DeleteShowToCategory,
		decodeDeleteShowToCategoryRequest,
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

func decodeGetShowsHomeRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
	id := chi.URLParam(r, "show_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	var req UpdateShowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	req.ID = id

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

	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}
	req.ShowID = showID

	return req, nil
}

func decodeUpdateEpisodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateEpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}
	req.ShowID = showID

	episodeID := chi.URLParam(r, "episode_id")
	if episodeID == "" {
		return nil, fmt.Errorf("%w: missed episodes id", ErrInvalidParameter)
	}
	req.ID = episodeID

	return req, nil
}

func decodeDeleteEpisodeByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	episodeID := chi.URLParam(r, "episode_id")
	if episodeID == "" {
		return nil, fmt.Errorf("%w: missed episode id", ErrInvalidParameter)
	}

	return DeleteEpisodeByIDRequest{
		ShowID:    showID,
		EpisodeID: episodeID,
	}, nil
}

func decodeGetEpisodeByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	episodeID := chi.URLParam(r, "episode_id")
	if episodeID == "" {
		return nil, fmt.Errorf("%w: missed episode id", ErrInvalidParameter)
	}

	return GetEpisodeByIDRequest{
		ShowID:    showID,
		EpisodeID: episodeID,
	}, nil
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

func decodeDeleteShowCategoryByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		return nil, fmt.Errorf("%w: missed category id", ErrInvalidParameter)
	}

	return categoryID, nil
}

func decodeGetShowCategoryByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		return nil, fmt.Errorf("%w: missed category id", ErrInvalidParameter)
	}

	return categoryID, nil
}

func decodeAddShowCategoriesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddShowsCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeUpdateShowCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateShowCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		return nil, fmt.Errorf("%w: missed category id", ErrInvalidParameter)
	}
	req.ID = categoryID

	return req, nil
}

func decodeDeleteShowToCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		return nil, fmt.Errorf("%w: missed category id", ErrInvalidParameter)
	}

	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	return DeleteShowToCategoryRequest{
		ShowID:     showID,
		CategoryID: categoryID,
	}, nil
}

func decodeDeleteShowToCategoryByShowIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	return showID, nil
}

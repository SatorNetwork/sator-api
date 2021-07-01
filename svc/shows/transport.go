package shows

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
	title             = "title"
	cover             = "cover"
	hasNewEpisode     = "has_new_episode"
	category          = "category"
	showId            = "show_id"
	episodeNumber     = "episode_number"
	description       = "description"
	releaseDate       = "release_date"
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

	r.Delete("/{show_id}/episodes", httptransport.NewServer(
		e.DeleteEpisodeByShowID,
		decodeDeleteEpisodeByShowIDRequest,
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
	b, err := strconv.ParseBool(r.URL.Query().Get(hasNewEpisode))
	if err != nil {
		return nil, fmt.Errorf("can not parse boolean value from string")
	}

	return AddShowRequest{
		Title:         r.URL.Query().Get(title),
		Cover:         r.URL.Query().Get(cover),
		HasNewEpisode: b,
		Category:      r.URL.Query().Get(category),
	}, nil
}

func decodeUpdateShowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := strconv.ParseBool(r.URL.Query().Get(hasNewEpisode))
	if err != nil {
		return nil, fmt.Errorf("can not parse boolean value from string")
	}

	return UpdateShowRequest{
		ID:            chi.URLParam(r, "show_id"),
		Title:         r.URL.Query().Get(title),
		Cover:         r.URL.Query().Get(cover),
		HasNewEpisode: b,
		Category:      r.URL.Query().Get(category),
	}, nil
}

func decodeDeleteShowByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "show_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed show_id", ErrInvalidParameter)
	}
	return id, nil
}

func decodeAddEpisodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return AddEpisodeRequest{
		ShowID:        r.URL.Query().Get(showId),
		EpisodeNumber: castStrToInt32(r.URL.Query().Get(episodeNumber)),
		Cover:         r.URL.Query().Get(cover),
		Title:         r.URL.Query().Get(title),
		Description:   r.URL.Query().Get(description),
		ReleaseDate:   r.URL.Query().Get(releaseDate),
	}, nil
}

func decodeUpdateEpisodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return UpdateEpisodeRequest{
		ID:            chi.URLParam(r, "id"),
		ShowID:        r.URL.Query().Get(showId),
		EpisodeNumber: castStrToInt32(r.URL.Query().Get(episodeNumber)),
		Cover:         r.URL.Query().Get(cover),
		Title:         r.URL.Query().Get(title),
		Description:   r.URL.Query().Get(description),
		ReleaseDate:   r.URL.Query().Get(releaseDate),
	}, nil
}

func decodeDeleteEpisodeByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed episode id", ErrInvalidParameter)
	}
	return id, nil
}

func decodeDeleteEpisodeByShowIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "show_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed show_id", ErrInvalidParameter)
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

package shows

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/lib/db"
	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/utils"

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

// Predefined request query keys
const (
	withNFT = "with_nft"
)

// MakeHTTPHandler ...
func MakeHTTPHandler(e Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	// shows
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

	// show categories
	r.Get("/categories", httptransport.NewServer(
		e.GetShowCategories,
		decodeGetShowCategoriesRequest,
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
		e.AddShowCategory,
		decodeAddShowCategoryRequest,
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

	// challenges
	r.Get("/{show_id}/challenges", httptransport.NewServer(
		e.GetShowChallenges,
		decodeGetShowChallengesRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	// Episodes
	r.Get("/{show_id}/episodes", httptransport.NewServer(
		e.GetEpisodesByShowID,
		decodeGetEpisodesByShowIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/episodes", httptransport.NewServer(
		e.GetActivatedUserEpisodes,
		decodeGetActivatedUserEpisodesRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/reviews", httptransport.NewServer(
		e.GetReviewsListByUserID,
		decodeGetReviewsListByUserIDRequest,
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

	r.Get("/filter/{category}", httptransport.NewServer(
		e.GetShowsByCategory,
		decodeGetShowsByCategoryRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{show_id}/episodes/{episode_id}/rate", httptransport.NewServer(
		e.RateEpisode,
		decodeRateEpisodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{show_id}/episodes/{episode_id}/reviews", httptransport.NewServer(
		e.ReviewEpisode,
		decodeReviewEpisodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{show_id}/episodes/{episode_id}/reviews", httptransport.NewServer(
		e.GetReviewsList,
		decodeGetReviewsListRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/reviews/{review_id}", httptransport.NewServer(
		e.DeleteReviewByID,
		decodeDeleteReviewByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/reviews/{review_id}/tips", httptransport.NewServer(
		e.SendTipsToReviewAuthor,
		decodeSendTipsToReviewAuthorRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	// Seasons
	r.Post("/{show_id}/seasons", httptransport.NewServer(
		e.AddSeason,
		decodeAddSeasonRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{show_id}/seasons/{season_id}", httptransport.NewServer(
		e.GetSeasonByID,
		decodeGetEpisodeByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{show_id}/seasons/{season_id}", httptransport.NewServer(
		e.DeleteSeasonByID,
		decodeDeleteSeasonByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{show_id}/claps", httptransport.NewServer(
		e.AddClapsForShow,
		decodeAddClapsForShowRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/reviews/{review_id}/{rating_type}", httptransport.NewServer(
		e.LikeDislikeEpisode,
		decodeLikeDislikeEpisodeRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrInvalidParameter) ||
		errors.Is(err, ErrAlreadyReviewed) ||
		errors.Is(err, ErrMaxClaps) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, ErrNotFound) || db.IsNotFoundError(err) {
		return http.StatusNotFound, err.Error()
	}

	if errors.Is(err, rbac.ErrAccessDenied) {
		return http.StatusForbidden, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeGetShowsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetShowsRequest{
		WithNFT: r.URL.Query().Get(withNFT),
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
		},
	}, nil
}

func decodeGetShowChallengesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetShowChallengesRequest{
		ShowID: chi.URLParam(r, "show_id"),
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
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
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
		},
	}, nil
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

	id := chi.URLParam(r, "show_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed show_id", ErrInvalidParameter)
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
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
		},
	}, nil
}

func decodeGetActivatedUserEpisodesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return utils.PaginationRequest{
		Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
		ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
	}, nil
}

func decodeGetReviewsListByUserIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return utils.PaginationRequest{
		Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
		ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
	}, nil
}

func decodeAddSeasonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddSeasonRequest
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

func decodeGetSeasonByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	seasonID := chi.URLParam(r, "season_id")
	if seasonID == "" {
		return nil, fmt.Errorf("%w: missed season id", ErrInvalidParameter)
	}

	return GetSeasonByIDRequest{
		ShowID:   showID,
		SeasonID: seasonID,
	}, nil
}

func decodeDeleteSeasonByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	seasonID := chi.URLParam(r, "season_id")
	if seasonID == "" {
		return nil, fmt.Errorf("%w: missed season id", ErrInvalidParameter)
	}

	return DeleteSeasonByIDRequest{
		ShowID:   showID,
		SeasonID: seasonID,
	}, nil
}

func decodeRateEpisodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RateEpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	episodeID := chi.URLParam(r, "episode_id")
	if episodeID == "" {
		return nil, fmt.Errorf("%w: missed episodes id", ErrInvalidParameter)
	}
	req.EpisodeID = episodeID

	return req, nil
}

func decodeReviewEpisodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ReviewEpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	episodeID := chi.URLParam(r, "episode_id")
	if episodeID == "" {
		return nil, fmt.Errorf("%w: missed episodes id", ErrInvalidParameter)
	}
	req.EpisodeID = episodeID

	return req, nil
}

func decodeGetReviewsListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetReviewsListRequest{
		EpisodeID: chi.URLParam(r, "episode_id"),
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
		},
	}, nil
}

func decodeDeleteReviewByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "review_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed review_id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeAddClapsForShowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	showID := chi.URLParam(r, "show_id")
	if showID == "" {
		return nil, fmt.Errorf("%w: missed show id", ErrInvalidParameter)
	}

	return showID, nil
}

func decodeLikeDislikeEpisodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req LikeDislikeEpisodeRequest
	id := chi.URLParam(r, "review_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed review id", ErrInvalidParameter)
	}
	param := chi.URLParam(r, "rating_type")
	if param == "" {
		return nil, fmt.Errorf("%w: missed like/dislike pamameter", ErrInvalidParameter)
	}
	req.ReviewID = id
	req.Param = param

	return req, nil
}

func decodeSendTipsToReviewAuthorRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SendTipsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	reviewID := chi.URLParam(r, "review_id")
	if reviewID == "" {
		return nil, fmt.Errorf("%w: missed review id", ErrInvalidParameter)
	}
	req.ReviewID = reviewID

	return req, nil
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

func decodeAddShowCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

func decodeGetShowCategoriesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetShowCategoriesRequest{
		WithDisabled: r.URL.Query().Get("with_disabled"),
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(utils.PageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(utils.ItemsPerPageParam)),
		},
	}, nil
}

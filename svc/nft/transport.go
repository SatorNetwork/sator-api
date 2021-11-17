package nft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/SatorNetwork/sator-api/internal/utils"

	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Predefined request query keys
const (
	pageParam         = "page"
	itemsPerPageParam = "items_per_page"
	relationId        = "relation_id"
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

	r.Post("/", httptransport.NewServer(
		e.CreateNFT,
		decodeCreateNFTRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/", httptransport.NewServer(
		e.GetNFTs,
		decodeGetNFTsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/filter/category/{category}", httptransport.NewServer(
		e.GetNFTsByCategory,
		decodeGetNFTsByCategoryRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/filter/show/{show_id}", httptransport.NewServer(
		e.GetNFTsByShowID,
		decodeGetNFTsByShowIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/filter/episode/{episode_id}", httptransport.NewServer(
		e.GetNFTsByEpisodeID,
		decodeGetNFTsByEpisodeIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/filter/user/{user_id}", httptransport.NewServer(
		e.GetNFTsByUserID,
		decodeGetNFTsByUserIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/{nft_id}", httptransport.NewServer(
		e.GetNFTByID,
		decodeGetNFTByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{nft_id}/buy", httptransport.NewServer(
		e.BuyNFT,
		decodeBuyNFTRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/{nft_id}", httptransport.NewServer(
		e.DeleteNFTItemByID,
		decodeDeleteNFTByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/{nft_id}", httptransport.NewServer(
		e.UpdateNFTItem,
		decodeUpdateNFTItemRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/categories", httptransport.NewServer(
		e.GetCategories,
		decodeGetCategoriesRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/home", httptransport.NewServer(
		e.GetMainScreenData,
		decodeGetMainScreenDataRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeCreateNFTRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req TransportNFT
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}
	return req, nil
}

func decodeGetNFTsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &GetNFTsWithFilterRequest{
		RelationID: r.URL.Query().Get(relationId),
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func decodeGetNFTsByCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	category := chi.URLParam(r, "category")
	if category == "" {
		return nil, fmt.Errorf("%w: missed category parameter", ErrInvalidParameter)
	}

	return &GetNFTsByCategoryRequest{
		Category: category,
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func decodeGetNFTsByShowIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	showId := chi.URLParam(r, "show_id")
	if showId == "" {
		return nil, fmt.Errorf("%w: missed category parameter", ErrInvalidParameter)
	}

	return &GetNFTsByShowIDRequest{
		ShowID: showId,
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func decodeGetNFTsByEpisodeIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	episodeId := chi.URLParam(r, "episode_id")
	if episodeId == "" {
		return nil, fmt.Errorf("%w: missed category parameter", ErrInvalidParameter)
	}

	return &GetNFTsByEpisodeIDRequest{
		EpisodeID: episodeId,
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func decodeGetNFTsByUserIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		return nil, fmt.Errorf("%w: missed user's ID parameter", ErrInvalidParameter)
	}

	return &GetNFTsByUserIDRequest{
		UserID: userId,
		PaginationRequest: utils.PaginationRequest{
			Page:         utils.StrToInt32(r.URL.Query().Get(pageParam)),
			ItemsPerPage: utils.StrToInt32(r.URL.Query().Get(itemsPerPageParam)),
		},
	}, nil
}

func decodeGetNFTByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	nftId := chi.URLParam(r, "nft_id")
	if nftId == "" {
		return nil, fmt.Errorf("%w: missed nft_id parameter", ErrInvalidParameter)
	}
	return nftId, nil
}

func decodeBuyNFTRequest(_ context.Context, r *http.Request) (interface{}, error) {
	nftId := chi.URLParam(r, "nft_id")
	if nftId == "" {
		return nil, fmt.Errorf("%w: missed nft_id parameter", ErrInvalidParameter)
	}
	return nftId, nil
}

func decodeGetCategoriesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return Empty{}, nil
}

func decodeGetMainScreenDataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return Empty{}, nil
}

func decodeDeleteNFTByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	nftID := chi.URLParam(r, "nft_id")
	if nftID == "" {
		return nil, fmt.Errorf("%w: missed nft id", ErrInvalidParameter)
	}

	return nftID, nil
}

func decodeUpdateNFTItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req UpdateNFTRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	nftID := chi.URLParam(r, "nft_id")
	if nftID == "" {
		return nil, fmt.Errorf("%w: missed nft id", ErrInvalidParameter)
	}

	id, err := uuid.Parse(nftID)
	if err != nil {
		return nil, fmt.Errorf("could not get nft id: %w", err)
	}

	req.ID = id

	return req, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

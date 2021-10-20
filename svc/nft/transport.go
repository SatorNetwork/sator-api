package nft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"

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

	r.Get("/filter/show_id/{show_id}/episode/{episode_id}", httptransport.NewServer(
		e.GetNFTsByShowID,
		decodeGetNFTsByShowIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/filter/user_id/{user_id}", httptransport.NewServer(
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
	return Empty{}, nil
}

func decodeGetNFTsByCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	category := chi.URLParam(r, "category")
	if category == "" {
		return nil, fmt.Errorf("%w: missed category parameter", ErrInvalidParameter)
	}

	return &GetNFTsByCategoryRequest{
		Category: category,
	}, nil
}

func decodeGetNFTsByShowIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	showId := chi.URLParam(r, "show_id")
	if showId == "" {
		return nil, fmt.Errorf("%w: missed category parameter", ErrInvalidParameter)
	}
	episodeId := chi.URLParam(r, "episode_id")
	if episodeId == "" {
		return nil, fmt.Errorf("%w: missed category parameter", ErrInvalidParameter)
	}

	return &GetNFTsByShowIDRequest{
		ShowID:    showId,
		EpisodeID: episodeId,
	}, nil
}

func decodeGetNFTsByUserIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		return nil, fmt.Errorf("%w: missed user's ID parameter", ErrInvalidParameter)
	}

	return &GetNFTsByUserIDRequest{
		UserID: userId,
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

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

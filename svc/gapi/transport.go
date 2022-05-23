package gapi

import (
	"context"
	"encoding/json"
	"errors"
	"log"
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

	r.Get("/get-status", httptransport.NewServer(
		e.GetStatus,
		decodeGetStatusRequest,
		customEncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/get-nft-packs", httptransport.NewServer(
		e.GetNFTPacks,
		decodeGetNFTPacksRequest,
		customEncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/buy-nft-pack", httptransport.NewServer(
		e.BuyNFTPack,
		decodeBuyNFTPackRequest,
		customEncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/craft-nft", httptransport.NewServer(
		e.CraftNFT,
		decodeCraftNFTRequest,
		customEncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/select-nft", httptransport.NewServer(
		e.SelectNFT,
		decodeSelectNFTRequest,
		customEncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/start-game", httptransport.NewServer(
		e.StartGame,
		decodeStartGameRequest,
		customEncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/finish-game", httptransport.NewServer(
		e.FinishGame,
		decodeFinishGameRequest,
		customEncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/claim-rewards", httptransport.NewServer(
		e.ClaimRewards,
		decodeClaimRewardsRequest,
		customEncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func decodeGetStatusRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetNFTPacksRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeBuyNFTPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req BuyNFTPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeCraftNFTRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CraftNFTRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeSelectNFTRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req SelectNFTRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeStartGameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req StartGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeFinishGameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req FinishGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeClaimRewardsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ClaimRewardsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrCouldNotVerifySignature) {
		log.Printf("could not verify signature: %v", err)
		return http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
	}

	if errors.Is(err, ErrNotAllNftsToCraftWereFound) {
		return http.StatusBadRequest, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

// customEncodeResponse extends the default EncodeResponse to sign the response
func customEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	signature, err := SignResponse([]byte("secret"), response)
	if err == nil && signature != "" {
		w.Header().Set("Signature", signature)
	}

	return httpencoder.EncodeResponse(ctx, w, response)
}

package gapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/SatorNetwork/sator-api/lib/jwt"
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
func MakeHTTPHandler(gameEndpoints Endpoints, nftPackEndpoints NFTPacksEndpoints, log logger, encodeResponse httptransport.EncodeResponseFunc) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(httpencoder.EncodeError(log, codeAndMessageFrom)),
		httptransport.ServerBefore(jwtkit.HTTPToContext()),
	}

	r.Get("/get-status", httptransport.NewServer(
		gameEndpoints.GetStatus,
		decodeGetStatusRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/get-nft-packs", httptransport.NewServer(
		gameEndpoints.GetNFTPacks,
		decodeGetNFTPacksRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/buy-nft-pack", httptransport.NewServer(
		gameEndpoints.BuyNFTPack,
		decodeBuyNFTPackRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/craft-nft", httptransport.NewServer(
		gameEndpoints.CraftNFT,
		decodeCraftNFTRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/select-nft", httptransport.NewServer(
		gameEndpoints.SelectNFT,
		decodeSelectNFTRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/start-game", httptransport.NewServer(
		gameEndpoints.StartGame,
		decodeStartGameRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/finish-game", httptransport.NewServer(
		gameEndpoints.FinishGame,
		decodeFinishGameRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/claim-rewards", httptransport.NewServer(
		gameEndpoints.ClaimRewards,
		decodeClaimRewardsRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	// Settings

	r.Get("/settings", httptransport.NewServer(
		gameEndpoints.GetSettings,
		decodeGetSettingsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/settings/value-types", httptransport.NewServer(
		gameEndpoints.GetSettingsValueTypes,
		decodeGetSettingsValueTypesRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/settings/{key}", httptransport.NewServer(
		gameEndpoints.GetSettingsByKey,
		decodeGetSettingsByKeyRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/settings", httptransport.NewServer(
		gameEndpoints.AddSetting,
		decodeAddSettingRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/settings/{key}", httptransport.NewServer(
		gameEndpoints.UpdateSetting,
		decodeUpdateSettingRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/settings/{key}", httptransport.NewServer(
		gameEndpoints.DeleteSetting,
		decodeDeleteSettingRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	// NFT packs
	r.Get("/nft-packs", httptransport.NewServer(
		nftPackEndpoints.GetNFTPacksListEndpoint,
		decodeGetNFTPacksRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/nft-packs", httptransport.NewServer(
		nftPackEndpoints.AddNFTPackEndpoint,
		decodeAddNFTPackRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/nft-packs/{id}", httptransport.NewServer(
		nftPackEndpoints.UpdateNFTPackEndpoint,
		decodeUpdateNFTPackRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/nft-packs/{id}", httptransport.NewServer(
		nftPackEndpoints.DeleteNFTPackEndpoint,
		decodeDeleteNFTPackRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/nft-packs/{id}/soft", httptransport.NewServer(
		nftPackEndpoints.SoftDeleteNFTPackEndpoint,
		decodeSoftDeleteNFTPackRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrCouldNotVerifySignature) {
		log.Printf("could not verify signature: %v", err)
		return http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
	}

	if errors.Is(err, ErrNotAllNftsToCraftWereFound) ||
		errors.Is(err, ErrNotEnoughNFTsToCraft) ||
		errors.Is(err, ErrNFTsToCraftHaveDifferentTypes) ||
		errors.Is(err, ErrNFTTypeLegendCannotBeCrafted) {
		return http.StatusBadRequest, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func EncodeResponseWithSignature(uv string) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		timestamp := time.Now().Unix()
		tokenString := ctx.Value(jwtkit.JWTTokenContextKey).(string)
		deviceID, _ := jwt.DeviceIDFromContext(ctx)
		signingKey := fmt.Sprintf("%s%s%d%s", uv, deviceID, timestamp, tokenString[len(tokenString)-3:])

		signature, err := SignResponse([]byte(signingKey), response)
		if err == nil && signature != "" {
			w.Header().Set("X-Signature", signature)
			w.Header().Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
		}

		return httpencoder.EncodeResponse(ctx, w, response)
	}
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

func decodeGetSettingsRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetSettingsValueTypesRequest(ctx context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeAddSettingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req AddGameSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateSettingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UpdateGameSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.Key = chi.URLParam(r, "key")

	return req, nil
}

func decodeDeleteSettingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return chi.URLParam(r, "key"), nil
}

func decodeAddNFTPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req AddNFTPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateNFTPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UpdateNFTPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.ID = chi.URLParam(r, "id")

	return req, nil
}

func decodeDeleteNFTPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return chi.URLParam(r, "id"), nil
}

func decodeSoftDeleteNFTPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return chi.URLParam(r, "id"), nil
}

func decodeGetSettingsByKeyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return chi.URLParam(r, "key"), nil
}

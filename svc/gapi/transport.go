package gapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/golang-jwt/jwt"
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

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

// EncodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func customEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response != nil {
		var (
			sig      string
			jsonBody []byte
			err      error
		)

		switch r := response.(type) {
		case httpencoder.Response, httpencoder.BoolResultResponse, httpencoder.ListResponse:
			jsonBody, err = json.Marshal(r)
			if err != nil {
				return ErrCouldNotSignResponse
			}
		case bool:
			jsonBody, err = json.Marshal(httpencoder.BoolResult(r))
			if err != nil {
				return ErrCouldNotSignResponse
			}
		default:
			jsonBody, err = json.Marshal(httpencoder.Response{Data: response})
			if err != nil {
				return ErrCouldNotSignResponse
			}
		}

		// Sign response body
		sig, err = jwt.SigningMethodHS256.Sign(string(jsonBody), []byte("secret"))
		if err != nil {
			return ErrCouldNotSignResponse
		}
		w.Header().Set("Signature", sig)
	}

	return httpencoder.EncodeResponse(ctx, w, response)
}

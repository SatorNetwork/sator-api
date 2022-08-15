package puzzle_game

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/SatorNetwork/sator-api/lib/httpencoder"
	"github.com/go-chi/chi/v5"
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

	r.Get("/admin/episode/{episode_id}", httptransport.NewServer(
		e.GetPuzzleGameByEpisodeID,
		decodeGetPuzzleGameByEpisodeIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/admin/{puzzle_game_id}", httptransport.NewServer(
		e.GetPuzzleGameByID,
		decodeGetPuzzleGameByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/admin", httptransport.NewServer(
		e.CreatePuzzleGame,
		decodeCreatePuzzleGameRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/admin/{puzzle_game_id}", httptransport.NewServer(
		e.UpdatePuzzleGame,
		decodeUpdatePuzzleGameRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/admin/{puzzle_game_id}/images/{image_id}", httptransport.NewServer(
		e.LinkImageToPuzzleGame,
		decodeLinkImageToPuzzleGameRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Delete("/admin/{puzzle_game_id}/images/{image_id}", httptransport.NewServer(
		e.DeleteImageByID,
		decodeDeleteImageByIDRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	// =========================================================================
	// User endpoints
	// =========================================================================

	r.Get("/episode/{episode_id}", httptransport.NewServer(
		e.GetPuzzleGameForUser,
		decodeGetPuzzleGameForUserRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/unlock-options", httptransport.NewServer(
		e.GetPuzzleGameUnlockOptions,
		decodeGetPuzzleGameUnlockOptionsRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{puzzle_game_id}/unlock", httptransport.NewServer(
		e.UnlockPuzzleGame,
		decodeUnlockPuzzleGameRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{puzzle_game_id}/start", httptransport.NewServer(
		e.StartPuzzleGame,
		decodeStartPuzzleGameRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/{puzzle_game_id}/tap-tile", httptransport.NewServer(
		e.TapTile,
		decodeTapTileRequest,
		httpencoder.EncodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden, err.Error()
	}

	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, err.Error()
	}

	if errors.Is(err, ErrInvalidParameter) {
		return http.StatusBadRequest, err.Error()
	}

	return httpencoder.CodeAndMessageFrom(err)
}

func decodeGetPuzzleGameByEpisodeIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "episode_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed episode id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeGetPuzzleGameByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "puzzle_game_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed puzzle game id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeCreatePuzzleGameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreatePuzzleGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	return req, nil
}

func decodeUpdatePuzzleGameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UpdatePuzzleGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	req.ID = chi.URLParam(r, "puzzle_game_id")
	if req.ID == "" {
		return nil, fmt.Errorf("%w: missed puzzle game id", ErrInvalidParameter)
	}

	return req, nil
}

func decodeLinkImageToPuzzleGameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	puzzleGameID := chi.URLParam(r, "puzzle_game_id")
	if puzzleGameID == "" {
		return nil, fmt.Errorf("%w: missed puzzle game id", ErrInvalidParameter)
	}

	imageID := chi.URLParam(r, "image_id")
	if imageID == "" {
		return nil, fmt.Errorf("%w: missed image id", ErrInvalidParameter)
	}

	return ImageToPuzzleGameRequest{
		PuzzleGameID: puzzleGameID,
		ImageID:      imageID,
	}, nil
}

func decodeDeleteImageByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	puzzleGameID := chi.URLParam(r, "puzzle_game_id")
	if puzzleGameID == "" {
		return nil, fmt.Errorf("%w: missed puzzle game id", ErrInvalidParameter)
	}

	imageID := chi.URLParam(r, "image_id")
	if imageID == "" {
		return nil, fmt.Errorf("%w: missed image id", ErrInvalidParameter)
	}

	return ImageToPuzzleGameRequest{
		PuzzleGameID: puzzleGameID,
		ImageID:      imageID,
	}, nil
}

func decodeGetPuzzleGameForUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "episode_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed episode id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeGetPuzzleGameUnlockOptionsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeUnlockPuzzleGameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UnlockPuzzleGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	id := chi.URLParam(r, "puzzle_game_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed puzzle game id", ErrInvalidParameter)
	}
	req.PuzzleGameID = id

	return req, nil
}

func decodeStartPuzzleGameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id := chi.URLParam(r, "puzzle_game_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed puzzle game id", ErrInvalidParameter)
	}

	return id, nil
}

func decodeTapTileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req TapTileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("could not decode request body: %w", err)
	}

	id := chi.URLParam(r, "puzzle_game_id")
	if id == "" {
		return nil, fmt.Errorf("%w: missed puzzle game id", ErrInvalidParameter)
	}
	req.PuzzleGameID = id

	return req, nil
}

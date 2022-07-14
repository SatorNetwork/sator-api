package puzzle_game

import (
	"context"
	"fmt"
	"log"

	"github.com/SatorNetwork/gopuzzlegame"

	"github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/validator"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetPuzzleGameByID        endpoint.Endpoint
		GetPuzzleGameByEpisodeID endpoint.Endpoint
		CreatePuzzleGame         endpoint.Endpoint
		UpdatePuzzleGame         endpoint.Endpoint

		LinkImageToPuzzleGame endpoint.Endpoint
		DeleteImageByID       endpoint.Endpoint

		GetPuzzleGameForUser endpoint.Endpoint
		UnlockPuzzleGame     endpoint.Endpoint
		StartPuzzleGame      endpoint.Endpoint

		GetPuzzleGameUnlockOptions endpoint.Endpoint
		TapTile                    endpoint.Endpoint
	}

	service interface {
		GetPuzzleGameByID(ctx context.Context, id uuid.UUID) (PuzzleGame, error)
		GetPuzzleGameByEpisodeID(ctx context.Context, episodeID uuid.UUID, isTestUser bool) (PuzzleGame, error)
		CreatePuzzleGame(ctx context.Context, epID uuid.UUID, prizePool float64, partsX int32) (PuzzleGame, error)
		UpdatePuzzleGame(ctx context.Context, id uuid.UUID, prizePool float64, partsX int32) (PuzzleGame, error)

		LinkImageToPuzzleGame(ctx context.Context, gameID, fileID uuid.UUID) error
		UnlinkImageFromPuzzleGame(ctx context.Context, gameID, fileID uuid.UUID) error

		GetPuzzleGameForUser(ctx context.Context, userID, episodeID uuid.UUID, isTestUser bool) (PuzzleGame, error)
		UnlockPuzzleGame(ctx context.Context, userID, puzzleGameID uuid.UUID, option string) (PuzzleGame, error)
		StartPuzzleGame(ctx context.Context, userID, puzzleGameID uuid.UUID) (PuzzleGame, error)

		GetPuzzleGameUnlockOptions(ctx context.Context) ([]PuzzleGameUnlockOption, error)
		TapTile(ctx context.Context, userID, puzzleGameID uuid.UUID, position gopuzzlegame.Position) (PuzzleGame, error)
	}

	CreatePuzzleGameRequest struct {
		EpisodeID string  `json:"episode_id" validate:"required,uuid"`
		PrizePool float64 `json:"prize_pool" validate:"min=0"`
		PartsX    int32   `json:"parts_x" validate:"required,min=3,max=10"`
	}

	UpdatePuzzleGameRequest struct {
		ID        string  `json:"episode_id" validate:"required,uuid"`
		PrizePool float64 `json:"prize_pool" validate:"required,min=0"`
		PartsX    int32   `json:"parts_x" validate:"required,min=3,max=10"`
	}

	ImageToPuzzleGameRequest struct {
		PuzzleGameID string `json:"puzzle_game_id" validate:"required,uuid"`
		ImageID      string `json:"image_id" validate:"required,uuid"`
	}

	UnlockPuzzleGameRequest struct {
		PuzzleGameID string `json:"puzzle_game_id" validate:"required,uuid"`
		UnlockOption string `json:"unlock_option" validate:"required"`
	}

	FinishPuzzleGameRequest struct {
		PuzzleGameID string `json:"puzzle_game_id" validate:"required,uuid"`
		Result       int    `json:"result" validate:"oneof=0 1 2"`
		StepsTaken   int    `json:"steps_taken" validate:"min=0"`
	}

	TapTileRequest struct {
		PuzzleGameID string `json:"puzzle_game_id" validate:"required,uuid"`
		X            int    `json:"x" validate:"min=0"`
		Y            int    `json:"y" validate:"min=0"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetPuzzleGameByID:        MakeGetPuzzleGameByIDEndpoint(s),
		GetPuzzleGameByEpisodeID: MakeGetPuzzleGameByEpisodeIDEndpoint(s),
		CreatePuzzleGame:         MakeCreatePuzzleGameEndpoint(s, validateFunc),
		UpdatePuzzleGame:         MakeUpdatePuzzleGameEndpoint(s, validateFunc),

		LinkImageToPuzzleGame: MakeLinkImageToPuzzleGameEndpoint(s, validateFunc),
		DeleteImageByID:       MakeDeleteImageByIDEndpoint(s, validateFunc),

		GetPuzzleGameForUser: MakeGetPuzzleGameForUserEndpoint(s),
		UnlockPuzzleGame:     MakeUnlockPuzzleGameEndpoint(s, validateFunc),
		StartPuzzleGame:      MakeStartPuzzleGameEndpoint(s),

		GetPuzzleGameUnlockOptions: MakeGetPuzzleGameUnlockOptionsEndpoint(s),
		TapTile:                    MakeTapTile(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetPuzzleGameByID = mdw(e.GetPuzzleGameByID)
			e.GetPuzzleGameByEpisodeID = mdw(e.GetPuzzleGameByEpisodeID)
			e.CreatePuzzleGame = mdw(e.CreatePuzzleGame)
			e.UpdatePuzzleGame = mdw(e.UpdatePuzzleGame)

			e.LinkImageToPuzzleGame = mdw(e.LinkImageToPuzzleGame)
			e.DeleteImageByID = mdw(e.DeleteImageByID)

			e.GetPuzzleGameForUser = mdw(e.GetPuzzleGameForUser)
			e.UnlockPuzzleGame = mdw(e.UnlockPuzzleGame)
			e.StartPuzzleGame = mdw(e.StartPuzzleGame)

			e.GetPuzzleGameUnlockOptions = mdw(e.GetPuzzleGameUnlockOptions)
			e.TapTile = mdw(e.TapTile)
		}
	}

	return e
}

func MakeGetPuzzleGameByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		res, err := s.GetPuzzleGameByID(ctx, uuid.MustParse(request.(string)))
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func MakeGetPuzzleGameByEpisodeIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		res, err := s.GetPuzzleGameByEpisodeID(ctx, uuid.MustParse(request.(string)), rbac.IsCurrentUserHasRole(ctx, rbac.RoleTestUser))
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func MakeCreatePuzzleGameEndpoint(s service, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(CreatePuzzleGameRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		res, err := s.CreatePuzzleGame(ctx, uuid.MustParse(req.EpisodeID), req.PrizePool, req.PartsX)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func MakeUpdatePuzzleGameEndpoint(s service, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(UpdatePuzzleGameRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		res, err := s.UpdatePuzzleGame(ctx, uuid.MustParse(req.ID), req.PrizePool, req.PartsX)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func MakeLinkImageToPuzzleGameEndpoint(s service, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(ImageToPuzzleGameRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		if err := s.LinkImageToPuzzleGame(
			ctx,
			uuid.MustParse(req.PuzzleGameID),
			uuid.MustParse(req.ImageID),
		); err != nil {
			return nil, err
		}

		return true, nil
	}
}

func MakeDeleteImageByIDEndpoint(s service, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(ImageToPuzzleGameRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		if err := s.UnlinkImageFromPuzzleGame(
			ctx,
			uuid.MustParse(req.PuzzleGameID),
			uuid.MustParse(req.ImageID),
		); err != nil {
			return false, err
		}

		return true, nil
	}
}

func MakeGetPuzzleGameForUserEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		res, err := s.GetPuzzleGameForUser(ctx, uid, uuid.MustParse(request.(string)), rbac.IsCurrentUserHasRole(ctx, rbac.RoleTestUser))
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func MakeUnlockPuzzleGameEndpoint(s service, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(UnlockPuzzleGameRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		res, err := s.UnlockPuzzleGame(ctx, uid, uuid.MustParse(req.PuzzleGameID), req.UnlockOption)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func MakeStartPuzzleGameEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		res, err := s.StartPuzzleGame(ctx, uid, uuid.MustParse(request.(string)))
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] StartPuzzleGame: %+v", res)

		return res, nil
	}
}

func MakeGetPuzzleGameUnlockOptionsEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		res, err := s.GetPuzzleGameUnlockOptions(ctx)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func MakeTapTile(s service, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		req := request.(TapTileRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		res, err := s.TapTile(ctx, uid, uuid.MustParse(req.PuzzleGameID), gopuzzlegame.Position{
			X: req.X,
			Y: req.Y,
		})
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

package puzzle_game

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	puzzle_game_repository "github.com/SatorNetwork/sator-api/svc/puzzle_game/repository"
)

type (
	Service struct {
		pgr puzzleGameRepository
	}

	puzzleGameRepository interface {
		GetPuzzleGameByID(ctx context.Context, id uuid.UUID) (puzzle_game_repository.PuzzleGame, error)
		CreatePuzzleGame(ctx context.Context, arg puzzle_game_repository.CreatePuzzleGameParams) (puzzle_game_repository.PuzzleGame, error)
		UpdatePuzzleGame(ctx context.Context, arg puzzle_game_repository.UpdatePuzzleGameParams) (puzzle_game_repository.PuzzleGame, error)
	}

	GetPuzzleGameByID struct {
		ID uuid.UUID `json:"id"`
	}

	CreatePuzzleGameRequest struct {
		EpisodeID uuid.UUID `json:"episode_id"`
		PrizePool float64   `json:"prize_pool"`
		PartsX    int32     `json:"parts_x"`
		PartsY    int32     `json:"parts_y"`
	}

	UpdatePuzzleGameRequest struct {
		EpisodeID uuid.UUID `json:"episode_id"`
		PrizePool float64   `json:"prize_pool"`
		PartsX    int32     `json:"parts_x"`
		PartsY    int32     `json:"parts_y"`
	}

	PuzzleGame struct {
		ID        uuid.UUID `json:"id"`
		EpisodeID uuid.UUID `json:"episode_id"`
		PrizePool float64   `json:"prize_pool"`
		PartsX    int32     `json:"parts_x"`
		PartsY    int32     `json:"parts_y"`
	}
)

func NewService(pgr puzzleGameRepository) *Service {
	s := &Service{
		pgr: pgr,
	}

	return s
}

//get puzzle game by id
//get puzzle game by episode_id (one item, not list)
//create puzzle game
//update puzzle game
//add an image to a puzzle game
//delete an image from a puzzle game
//get image list by puzzle game id

func (s *Service) GetPuzzleGameByID(ctx context.Context, req *GetPuzzleGameByID) (*PuzzleGame, error) {
	puzzleGame, err := s.pgr.GetPuzzleGameByID(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get puzzle game by id")
	}

	return NewPuzzleGameFromSQLC(&puzzleGame), nil
}

func (s *Service) CreatePuzzleGame(ctx context.Context, req *CreatePuzzleGameRequest) (*PuzzleGame, error) {
	puzzleGame, err := s.pgr.CreatePuzzleGame(ctx, puzzle_game_repository.CreatePuzzleGameParams{
		EpisodeID: req.EpisodeID,
		PrizePool: req.PrizePool,
		PartsX:    req.PartsX,
		PartsY:    req.PartsY,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't create puzzle game")
	}

	return NewPuzzleGameFromSQLC(&puzzleGame), nil
}

func (s *Service) UpdatePuzzleGame(ctx context.Context, req *UpdatePuzzleGameRequest) (*PuzzleGame, error) {
	puzzleGame, err := s.pgr.UpdatePuzzleGame(ctx, puzzle_game_repository.UpdatePuzzleGameParams{
		EpisodeID: req.EpisodeID,
		PrizePool: req.PrizePool,
		PartsX:    req.PartsX,
		PartsY:    req.PartsY,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't update puzzle game")
	}

	return NewPuzzleGameFromSQLC(&puzzleGame), nil
}

//func (s *Service) DeleteLink(ctx context.Context, id uuid.UUID) error {
//	if err := s.tpr.DeleteTradingPlatformLink(ctx, id); err != nil {
//		return errors.Wrap(err, "can't delete trading platform link")
//	}
//
//	return nil
//}
//
//func (s *Service) GetLinks(ctx context.Context, req *utils.PaginationRequest) ([]*Link, error) {
//	links, err := s.tpr.GetTradingPlatformLinks(ctx, trading_platforms_repository.GetTradingPlatformLinksParams{
//		Limit:  req.Limit(),
//		Offset: req.Offset(),
//	})
//	if err != nil {
//		return nil, errors.Wrap(err, "can't get trading platform links")
//	}
//
//	return NewLinksFromSQLC(links), nil
//}

func NewPuzzleGameFromSQLC(pg *puzzle_game_repository.PuzzleGame) *PuzzleGame {
	return &PuzzleGame{
		ID:        pg.ID,
		EpisodeID: pg.EpisodeID,
		PrizePool: pg.PrizePool,
		PartsX:    pg.PartsX,
		PartsY:    pg.PartsY,
	}
}

//func NewLinksFromSQLC(sqlcLinks []trading_platforms_repository.TradingPlatformLink) []*Link {
//	links := make([]*Link, 0, len(sqlcLinks))
//	for _, sqlcLink := range sqlcLinks {
//		links = append(links, NewLinkFromSQLC(&sqlcLink))
//	}
//
//	return links
//}

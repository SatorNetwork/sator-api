package puzzle_game

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	files_svc "github.com/SatorNetwork/sator-api/svc/files"
	puzzle_game_repository "github.com/SatorNetwork/sator-api/svc/puzzle_game/repository"
)

type (
	Service struct {
		pgr      puzzleGameRepository
		filesSvc filesService
	}

	puzzleGameRepository interface {
		GetPuzzleGameByID(ctx context.Context, id uuid.UUID) (puzzle_game_repository.PuzzleGame, error)
		CreatePuzzleGame(ctx context.Context, arg puzzle_game_repository.CreatePuzzleGameParams) (puzzle_game_repository.PuzzleGame, error)
		UpdatePuzzleGame(ctx context.Context, arg puzzle_game_repository.UpdatePuzzleGameParams) (puzzle_game_repository.PuzzleGame, error)
		AddImageToPuzzleGame(ctx context.Context, arg puzzle_game_repository.AddImageToPuzzleGameParams) error
		DeleteImageFromPuzzleGame(ctx context.Context, arg puzzle_game_repository.DeleteImageFromPuzzleGameParams) error
	}

	filesService interface {
		AddFile(ctx context.Context, it files_svc.File, file []byte, fileHeader *multipart.FileHeader) (files_svc.File, error)
	}

	Empty struct{}

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

	AddImageToPuzzleGameRequest struct {
		RawImage     string    `json:"raw_image"`
		Filename string `json:"filename"`
		FileHeader *multipart.FileHeader
		PuzzleGameID uuid.UUID `json:"puzzle_game_id"`
	}

	DeleteImageFromPuzzleGameRequest struct {
		FileID       uuid.UUID `json:"file_id"`
		PuzzleGameID uuid.UUID `json:"puzzle_game_id"`
	}
)

func NewService(pgr puzzleGameRepository, filesSvc filesService) *Service {
	s := &Service{
		pgr:      pgr,
		filesSvc: filesSvc,
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

//get puzzle game by episode_id (one item, not list)

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

//add an image to a puzzle game
//delete an image from a puzzle game
//get image list by puzzle game id
func (s *Service) AddImageToPuzzleGame(ctx context.Context, req *AddImageToPuzzleGameRequest) (*Empty, error) {
	var fileID uuid.UUID
	// TODO:
	// ctx context.Context, it files_svc.File, file []byte, fileHeader *multipart.FileHeader
	s.filesSvc.AddFile(ctx, files_svc.File{Filename: req.Filename}, []byte(req.RawImage), req.FileHeader)

	err := s.pgr.AddImageToPuzzleGame(ctx, puzzle_game_repository.AddImageToPuzzleGameParams{
		FileID:       fileID,
		PuzzleGameID: req.PuzzleGameID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't add image to puzzle game")
	}

	return &Empty{}, nil
}

func (s *Service) DeleteImageFromPuzzleGame(ctx context.Context, req *DeleteImageFromPuzzleGameRequest) (*Empty, error) {
	err := s.pgr.DeleteImageFromPuzzleGame(ctx, puzzle_game_repository.DeleteImageFromPuzzleGameParams{
		FileID:       req.FileID,
		PuzzleGameID: req.PuzzleGameID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't delete image from puzzle game")
	}

	return &Empty{}, nil
}

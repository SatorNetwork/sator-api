package puzzle_game

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	filesRepo "github.com/SatorNetwork/sator-api/svc/files/repository"
	puzzleGameRepo "github.com/SatorNetwork/sator-api/svc/puzzle_game/repository"
	showsRepo "github.com/SatorNetwork/sator-api/svc/shows/repository"
)

type DB struct {
	dbClient             *sql.DB
	showsRepository      *showsRepo.Queries
	filesRepository      *filesRepo.Queries
	puzzleGameRepository *puzzleGameRepo.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	puzzleGameRepository, err := puzzleGameRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "puzzleGameRepository error")
	}

	showsRepository, err := showsRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "showsRepository error")
	}

	filesRepository, err := filesRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "filesRepository error")
	}

	return &DB{
		puzzleGameRepository: puzzleGameRepository,
		dbClient:             dbClient,
		showsRepository:      showsRepository,
		filesRepository:      filesRepository,
	}, nil
}

func (db *DB) GetPuzzleGameByEpisodeID(ctx context.Context, episodeID uuid.UUID) (*puzzleGameRepo.PuzzleGame, error) {
	pg, err := db.puzzleGameRepository.GetPuzzleGameByEpisodeID(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("could not get puzzle game by episodeID: %v: %w", episodeID, err)
	}

	return &pg, nil
}

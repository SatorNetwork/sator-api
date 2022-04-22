package shows

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	shows_repo "github.com/SatorNetwork/sator-api/svc/shows/repository"
)

type DB struct {
	dbClient        *sql.DB
	showsRepository *shows_repo.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	showsRepository, err := shows_repo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare shows repository")
	}

	return &DB{
		dbClient:        dbClient,
		showsRepository: showsRepository,
	}, nil
}

func (db *DB) GetShowsByTitle(ctx context.Context, title string) ([]shows_repo.Show, error) {
	shows, err := db.showsRepository.GetShowsByTitle(ctx, title)
	if err != nil {
		return nil, fmt.Errorf("could not get shows by title: %v: %w", title, err)
	}

	return shows, nil
}

func (db *DB) GetEpisodesIDByShowID(showID uuid.UUID) ([]uuid.UUID, error) {
	var result []uuid.UUID
	rows, err := db.dbClient.Query("SELECT id FROM public.episodes WHERE show_id = $1", showID)
	if err != nil {
		return nil, fmt.Errorf("could not get episodes by showID: %v: %w", showID, err)
	}

	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}

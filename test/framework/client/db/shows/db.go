package shows

import (
	"database/sql"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	shows_repo "github.com/SatorNetwork/sator-api/svc/shows/repository"
)

type DB struct {
	showsRepository *shows_repo.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	showsRepository, err := shows_repo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare shows repository")
	}

	return &DB{
		showsRepository: showsRepository,
	}, nil
}

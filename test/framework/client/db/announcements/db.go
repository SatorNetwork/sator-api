package announcements

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	announcement_repository "github.com/SatorNetwork/sator-api/svc/announcement/repository"
)

type DB struct {
	announcementsRepository *announcement_repository.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	announcementsRepository, err := announcement_repository.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare announcement repository")
	}

	return &DB{
		announcementsRepository: announcementsRepository,
	}, nil
}

func (db *DB) Repository() *announcement_repository.Queries {
	return db.announcementsRepository
}

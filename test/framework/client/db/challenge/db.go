package challenge

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
	shows_repository "github.com/SatorNetwork/sator-api/svc/shows/repository"
)

type DB struct {
	challengeRepository *challengeRepo.Queries
	showsRepository     *shows_repository.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	challengeRepository, err := challengeRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "challengeRepo error")
	}
	showsRepository, err := shows_repository.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare shows repository")
	}

	return &DB{
		challengeRepository: challengeRepository,
		showsRepository:     showsRepository,
	}, nil
}

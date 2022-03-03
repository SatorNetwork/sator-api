package challenge

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	challengeRepo "github.com/SatorNetwork/sator-api/svc/challenge/repository"
)

type DB struct {
	challengeRepository *challengeRepo.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	challengeRepository, err := challengeRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "challengeRepo error")
	}

	return &DB{
		challengeRepository: challengeRepository,
	}, nil
}

package iap

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	iap_repository "github.com/SatorNetwork/sator-api/svc/iap/repository"
)

type DB struct {
	iapRepository *iap_repository.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	iapRepository, err := iap_repository.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare iap repository")
	}

	return &DB{
		iapRepository: iapRepository,
	}, nil
}

package wallet

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	walletRepo "github.com/SatorNetwork/sator-api/svc/wallet/repository"
)

type DB struct {
	walletRepository *walletRepo.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	walletRepository, err := walletRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare wallet repository")
	}

	return &DB{
		walletRepository: walletRepository,
	}, nil
}

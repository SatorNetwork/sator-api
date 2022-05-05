package iap

import (
	"context"

	"github.com/pkg/errors"

	iap_repository "github.com/SatorNetwork/sator-api/svc/iap/repository"
)

func (db *DB) Bootstrap(ctx context.Context) error {
	_, err := db.iapRepository.CreateIapProduct(ctx, iap_repository.CreateIapProductParams{
		ID:         "test2",
		PriceInSao: 100,
	})
	if err != nil {
		return errors.Wrap(err, "can't create iap product")
	}

	return nil
}

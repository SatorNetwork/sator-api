package wallet

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
	"github.com/google/uuid"
)

func (db *DB) Bootstrap(ctx context.Context) error {
	_, err := db.walletRepository.AddStakeLevel(ctx, repository.AddStakeLevelParams{
		MinStakeAmount: sql.NullFloat64{
			Float64: 1,
			Valid:   true,
		},
		MinDaysAmount: sql.NullInt32{
			Int32: 1,
			Valid: true,
		},
		Title:    "mock-title",
		Subtitle: "mock-subtitle",
		Multiplier: sql.NullInt32{
			Int32: 1,
			Valid: true,
		},
		Disabled: sql.NullBool{
			Bool:  false,
			Valid: true,
		},
	})
	if err != nil {
		return errors.Wrap(err, "can't AddStakeLevel")
	}

	return nil
}

func (db *DB) SetEmptyStake(userID uuid.UUID) error {
	wallet, err := db.walletRepository.GetWalletByUserIDAndType(context.Background(), repository.GetWalletByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: "sao",
	})
	if err != nil {
		return err
	}

	_, err = db.walletRepository.AddStake(context.Background(), repository.AddStakeParams{
		UserID:      userID,
		WalletID:    wallet.ID,
		StakeAmount: 1,
	})
	if err != nil {
		return err
	}

	return nil
}

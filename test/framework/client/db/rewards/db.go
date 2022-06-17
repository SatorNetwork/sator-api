package rewards

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	rewardsRepo "github.com/SatorNetwork/sator-api/svc/rewards/repository"
)

type DB struct {
	authRepository *rewardsRepo.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	authRepository, err := rewardsRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "authRepository error")
	}

	return &DB{
		authRepository: authRepository,
	}, nil
}

func (db *DB) DepositRewards(ctx context.Context, userID uuid.UUID, amount float64) error {
	if err := db.authRepository.AddTransaction(ctx, rewardsRepo.AddTransactionParams{
		UserID:          userID,
		TransactionType: 1,
		Amount:          amount,
	}); err != nil {
		return fmt.Errorf("error to deposit rewards for user: %v: %w", userID, err)
	}

	return nil
}

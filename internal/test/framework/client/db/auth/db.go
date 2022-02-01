package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	authRepo "github.com/SatorNetwork/sator-api/svc/auth/repository"
)

type DB struct {
	authRepository *authRepo.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	authRepository, err := authRepo.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "authRepository error")
	}

	return &DB{
		authRepository: authRepository,
	}, nil
}

func (db *DB) UpdateKYCStatus(ctx context.Context, email, status string) error {
	u, err := db.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("could not get user by email %s", email)
	}

	if err := db.authRepository.UpdateKYCStatus(ctx, authRepo.UpdateKYCStatusParams{
		KycStatus: status,
		ID:        u.ID,
	}); err != nil {
		return fmt.Errorf("could not update kyc status for user: %v: %w", u.ID, err)
	}

	return nil
}

func (db *DB) GetUserIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	u, err := db.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get user by email %s", email)
	}

	return u.ID, nil
}

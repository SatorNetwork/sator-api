package firebase

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	firebase_repository "github.com/SatorNetwork/sator-api/svc/firebase/repository"
)

type DB struct {
	firebaseRepository *firebase_repository.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	firebaseRepository, err := firebase_repository.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare firebase repository")
	}

	return &DB{
		firebaseRepository: firebaseRepository,
	}, nil
}

func (db *DB) Repository() *firebase_repository.Queries {
	return db.firebaseRepository
}

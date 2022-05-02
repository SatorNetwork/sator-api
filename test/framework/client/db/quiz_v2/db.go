package quiz_v2

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	quiz_v2_repository "github.com/SatorNetwork/sator-api/svc/quiz_v2/repository"
)

type DB struct {
	quizV2Repository *quiz_v2_repository.Queries
}

func New(dbClient *sql.DB) (*DB, error) {
	ctx := context.Background()
	quizV2Repository, err := quiz_v2_repository.Prepare(ctx, dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare quiz v2 repository")
	}

	return &DB{
		quizV2Repository: quizV2Repository,
	}, nil
}

func (db *DB) Repository() *quiz_v2_repository.Queries {
	return db.quizV2Repository
}

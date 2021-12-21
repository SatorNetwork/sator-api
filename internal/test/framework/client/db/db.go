package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/internal/test/framework/client/db/auth"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/db/challenge"
)

type DB struct {
	dbClient *sql.DB

	challengeDB *challenge.DB
	authDB      *auth.DB
}

func New() (*DB, error) {
	// TODO: use env var instead
	dbConnString := "postgresql://pguser:pgpass@127.0.0.1/pgdb?sslmode=disable"
	dbClient, err := sql.Open("postgres", dbConnString)
	if err != nil {
		return nil, errors.Wrap(err, "init db connection error")
	}

	authDB, err := auth.New(dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't create auth db")
	}

	challengeDB, err := challenge.New(dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't create challenge db")
	}

	return &DB{
		dbClient:    dbClient,
		challengeDB: challengeDB,
		authDB:      authDB,
	}, nil
}

func (db *DB) Bootstrap(ctx context.Context) error {
	if err := db.challengeDB.Bootstrap(ctx); err != nil {
		return err
	}

	return nil
}

func (db *DB) ChallengeDB() *challenge.DB {
	return db.challengeDB
}

func (db *DB) AuthDB() *auth.DB {
	return db.authDB
}

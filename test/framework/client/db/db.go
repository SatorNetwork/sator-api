package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/test/framework/client/db/auth"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/challenge"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/shows"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/wallet"
)

type DB struct {
	dbClient *sql.DB

	challengeDB *challenge.DB
	walletDB    *wallet.DB
	authDB      *auth.DB
	showsDB     *shows.DB
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

	walletDB, err := wallet.New(dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't create wallet db")
	}

	showsDB, err := shows.New(dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't create shows db")
	}

	return &DB{
		dbClient:    dbClient,
		challengeDB: challengeDB,
		authDB:      authDB,
		walletDB:    walletDB,
		showsDB:     showsDB,
	}, nil
}

func (db *DB) Bootstrap(ctx context.Context) error {
	if err := db.challengeDB.Bootstrap(ctx); err != nil {
		return errors.Wrap(err, "can't bootstrap challenge db")
	}

	if err := db.walletDB.Bootstrap(ctx); err != nil {
		return errors.Wrap(err, "can't bootstrap wallet db")
	}

	if err := db.showsDB.Bootstrap(ctx); err != nil {
		return errors.Wrap(err, "can't bootstrap shows db")
	}

	return nil
}

func (db *DB) ChallengeDB() *challenge.DB {
	return db.challengeDB
}

func (db *DB) AuthDB() *auth.DB {
	return db.authDB
}

func (db *DB) WalletDB() *wallet.DB {
	return db.walletDB
}
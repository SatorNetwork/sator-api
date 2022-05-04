package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/test/framework/client/db/auth"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/challenge"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/puzzle_game"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/quiz_v2"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/rewards"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/shows"
	"github.com/SatorNetwork/sator-api/test/framework/client/db/wallet"
)

type DB struct {
	dbClient *sql.DB

	challengeDB  *challenge.DB
	walletDB     *wallet.DB
	authDB       *auth.DB
	showsDB      *shows.DB
	quizV2DB     *quiz_v2.DB
	puzzleGameDB *puzzle_game.DB
	rewardsDB	*rewards.DB
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

	quizV2DB, err := quiz_v2.New(dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't create quiz v2 db")
	}

	puzzleGameDB, err := puzzle_game.New(dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't create puzzle game db")
	}

	rewardsDB, err := rewards.New(dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't create rewards db")
	}

	return &DB{
		dbClient:     dbClient,
		challengeDB:  challengeDB,
		authDB:       authDB,
		walletDB:     walletDB,
		showsDB:      showsDB,
		quizV2DB:     quizV2DB,
		puzzleGameDB: puzzleGameDB,
		rewardsDB:   rewardsDB,
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

	if err := db.puzzleGameDB.Bootstrap(ctx); err != nil {
		return errors.Wrap(err, "can't bootstrap puzzle game db")
	}

	return nil
}

func (db *DB) Client() *sql.DB {
	return db.dbClient
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

func (db *DB) ShowsDB() *shows.DB {
	return db.showsDB
}

func (db *DB) QuizV2DB() *quiz_v2.DB {
	return db.quizV2DB
}

func (db *DB) PuzzleGameDB() *puzzle_game.DB {
	return db.puzzleGameDB
}

func (db *DB) RewardsDB() *rewards.DB {
	return db.rewardsDB
}
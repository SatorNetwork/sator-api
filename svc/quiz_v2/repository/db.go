// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.cleanUpStmt, err = db.PrepareContext(ctx, cleanUp); err != nil {
		return nil, fmt.Errorf("error preparing query CleanUp: %w", err)
	}
	if q.countPlayersInRoomStmt, err = db.PrepareContext(ctx, countPlayersInRoom); err != nil {
		return nil, fmt.Errorf("error preparing query CountPlayersInRoom: %w", err)
	}
	if q.getDistributedRewardsByChallengeIDStmt, err = db.PrepareContext(ctx, getDistributedRewardsByChallengeID); err != nil {
		return nil, fmt.Errorf("error preparing query GetDistributedRewardsByChallengeID: %w", err)
	}
	if q.registerNewPlayerStmt, err = db.PrepareContext(ctx, registerNewPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query RegisterNewPlayer: %w", err)
	}
	if q.registerNewQuizStmt, err = db.PrepareContext(ctx, registerNewQuiz); err != nil {
		return nil, fmt.Errorf("error preparing query RegisterNewQuiz: %w", err)
	}
	if q.unregisterPlayerStmt, err = db.PrepareContext(ctx, unregisterPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query UnregisterPlayer: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.cleanUpStmt != nil {
		if cerr := q.cleanUpStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing cleanUpStmt: %w", cerr)
		}
	}
	if q.countPlayersInRoomStmt != nil {
		if cerr := q.countPlayersInRoomStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countPlayersInRoomStmt: %w", cerr)
		}
	}
	if q.getDistributedRewardsByChallengeIDStmt != nil {
		if cerr := q.getDistributedRewardsByChallengeIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDistributedRewardsByChallengeIDStmt: %w", cerr)
		}
	}
	if q.registerNewPlayerStmt != nil {
		if cerr := q.registerNewPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing registerNewPlayerStmt: %w", cerr)
		}
	}
	if q.registerNewQuizStmt != nil {
		if cerr := q.registerNewQuizStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing registerNewQuizStmt: %w", cerr)
		}
	}
	if q.unregisterPlayerStmt != nil {
		if cerr := q.unregisterPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing unregisterPlayerStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                     DBTX
	tx                                     *sql.Tx
	cleanUpStmt                            *sql.Stmt
	countPlayersInRoomStmt                 *sql.Stmt
	getDistributedRewardsByChallengeIDStmt *sql.Stmt
	registerNewPlayerStmt                  *sql.Stmt
	registerNewQuizStmt                    *sql.Stmt
	unregisterPlayerStmt                   *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                     tx,
		tx:                                     tx,
		cleanUpStmt:                            q.cleanUpStmt,
		countPlayersInRoomStmt:                 q.countPlayersInRoomStmt,
		getDistributedRewardsByChallengeIDStmt: q.getDistributedRewardsByChallengeIDStmt,
		registerNewPlayerStmt:                  q.registerNewPlayerStmt,
		registerNewQuizStmt:                    q.registerNewQuizStmt,
		unregisterPlayerStmt:                   q.unregisterPlayerStmt,
	}
}

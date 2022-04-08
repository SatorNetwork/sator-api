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
	if q.createPuzzleGameStmt, err = db.PrepareContext(ctx, createPuzzleGame); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePuzzleGame: %w", err)
	}
	if q.finishPuzzleGameStmt, err = db.PrepareContext(ctx, finishPuzzleGame); err != nil {
		return nil, fmt.Errorf("error preparing query FinishPuzzleGame: %w", err)
	}
	if q.getPuzzleGameByEpisodeIDStmt, err = db.PrepareContext(ctx, getPuzzleGameByEpisodeID); err != nil {
		return nil, fmt.Errorf("error preparing query GetPuzzleGameByEpisodeID: %w", err)
	}
	if q.getPuzzleGameByIDStmt, err = db.PrepareContext(ctx, getPuzzleGameByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetPuzzleGameByID: %w", err)
	}
	if q.getPuzzleGameCurrentAttemptStmt, err = db.PrepareContext(ctx, getPuzzleGameCurrentAttempt); err != nil {
		return nil, fmt.Errorf("error preparing query GetPuzzleGameCurrentAttempt: %w", err)
	}
	if q.getPuzzleGameImageIDsStmt, err = db.PrepareContext(ctx, getPuzzleGameImageIDs); err != nil {
		return nil, fmt.Errorf("error preparing query GetPuzzleGameImageIDs: %w", err)
	}
	if q.getPuzzleGameUnlockOptionStmt, err = db.PrepareContext(ctx, getPuzzleGameUnlockOption); err != nil {
		return nil, fmt.Errorf("error preparing query GetPuzzleGameUnlockOption: %w", err)
	}
	if q.getPuzzleGameUnlockOptionsStmt, err = db.PrepareContext(ctx, getPuzzleGameUnlockOptions); err != nil {
		return nil, fmt.Errorf("error preparing query GetPuzzleGameUnlockOptions: %w", err)
	}
	if q.getUserAvailableStepsStmt, err = db.PrepareContext(ctx, getUserAvailableSteps); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserAvailableSteps: %w", err)
	}
	if q.linkImageToPuzzleGameStmt, err = db.PrepareContext(ctx, linkImageToPuzzleGame); err != nil {
		return nil, fmt.Errorf("error preparing query LinkImageToPuzzleGame: %w", err)
	}
	if q.startPuzzleGameStmt, err = db.PrepareContext(ctx, startPuzzleGame); err != nil {
		return nil, fmt.Errorf("error preparing query StartPuzzleGame: %w", err)
	}
	if q.unlinkImageFromPuzzleGameStmt, err = db.PrepareContext(ctx, unlinkImageFromPuzzleGame); err != nil {
		return nil, fmt.Errorf("error preparing query UnlinkImageFromPuzzleGame: %w", err)
	}
	if q.unlockPuzzleGameStmt, err = db.PrepareContext(ctx, unlockPuzzleGame); err != nil {
		return nil, fmt.Errorf("error preparing query UnlockPuzzleGame: %w", err)
	}
	if q.updatePuzzleGameStmt, err = db.PrepareContext(ctx, updatePuzzleGame); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePuzzleGame: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createPuzzleGameStmt != nil {
		if cerr := q.createPuzzleGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPuzzleGameStmt: %w", cerr)
		}
	}
	if q.finishPuzzleGameStmt != nil {
		if cerr := q.finishPuzzleGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing finishPuzzleGameStmt: %w", cerr)
		}
	}
	if q.getPuzzleGameByEpisodeIDStmt != nil {
		if cerr := q.getPuzzleGameByEpisodeIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPuzzleGameByEpisodeIDStmt: %w", cerr)
		}
	}
	if q.getPuzzleGameByIDStmt != nil {
		if cerr := q.getPuzzleGameByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPuzzleGameByIDStmt: %w", cerr)
		}
	}
	if q.getPuzzleGameCurrentAttemptStmt != nil {
		if cerr := q.getPuzzleGameCurrentAttemptStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPuzzleGameCurrentAttemptStmt: %w", cerr)
		}
	}
	if q.getPuzzleGameImageIDsStmt != nil {
		if cerr := q.getPuzzleGameImageIDsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPuzzleGameImageIDsStmt: %w", cerr)
		}
	}
	if q.getPuzzleGameUnlockOptionStmt != nil {
		if cerr := q.getPuzzleGameUnlockOptionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPuzzleGameUnlockOptionStmt: %w", cerr)
		}
	}
	if q.getPuzzleGameUnlockOptionsStmt != nil {
		if cerr := q.getPuzzleGameUnlockOptionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPuzzleGameUnlockOptionsStmt: %w", cerr)
		}
	}
	if q.getUserAvailableStepsStmt != nil {
		if cerr := q.getUserAvailableStepsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserAvailableStepsStmt: %w", cerr)
		}
	}
	if q.linkImageToPuzzleGameStmt != nil {
		if cerr := q.linkImageToPuzzleGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing linkImageToPuzzleGameStmt: %w", cerr)
		}
	}
	if q.startPuzzleGameStmt != nil {
		if cerr := q.startPuzzleGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing startPuzzleGameStmt: %w", cerr)
		}
	}
	if q.unlinkImageFromPuzzleGameStmt != nil {
		if cerr := q.unlinkImageFromPuzzleGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing unlinkImageFromPuzzleGameStmt: %w", cerr)
		}
	}
	if q.unlockPuzzleGameStmt != nil {
		if cerr := q.unlockPuzzleGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing unlockPuzzleGameStmt: %w", cerr)
		}
	}
	if q.updatePuzzleGameStmt != nil {
		if cerr := q.updatePuzzleGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePuzzleGameStmt: %w", cerr)
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
	db                              DBTX
	tx                              *sql.Tx
	createPuzzleGameStmt            *sql.Stmt
	finishPuzzleGameStmt            *sql.Stmt
	getPuzzleGameByEpisodeIDStmt    *sql.Stmt
	getPuzzleGameByIDStmt           *sql.Stmt
	getPuzzleGameCurrentAttemptStmt *sql.Stmt
	getPuzzleGameImageIDsStmt       *sql.Stmt
	getPuzzleGameUnlockOptionStmt   *sql.Stmt
	getPuzzleGameUnlockOptionsStmt  *sql.Stmt
	getUserAvailableStepsStmt       *sql.Stmt
	linkImageToPuzzleGameStmt       *sql.Stmt
	startPuzzleGameStmt             *sql.Stmt
	unlinkImageFromPuzzleGameStmt   *sql.Stmt
	unlockPuzzleGameStmt            *sql.Stmt
	updatePuzzleGameStmt            *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                              tx,
		tx:                              tx,
		createPuzzleGameStmt:            q.createPuzzleGameStmt,
		finishPuzzleGameStmt:            q.finishPuzzleGameStmt,
		getPuzzleGameByEpisodeIDStmt:    q.getPuzzleGameByEpisodeIDStmt,
		getPuzzleGameByIDStmt:           q.getPuzzleGameByIDStmt,
		getPuzzleGameCurrentAttemptStmt: q.getPuzzleGameCurrentAttemptStmt,
		getPuzzleGameImageIDsStmt:       q.getPuzzleGameImageIDsStmt,
		getPuzzleGameUnlockOptionStmt:   q.getPuzzleGameUnlockOptionStmt,
		getPuzzleGameUnlockOptionsStmt:  q.getPuzzleGameUnlockOptionsStmt,
		getUserAvailableStepsStmt:       q.getUserAvailableStepsStmt,
		linkImageToPuzzleGameStmt:       q.linkImageToPuzzleGameStmt,
		startPuzzleGameStmt:             q.startPuzzleGameStmt,
		unlinkImageFromPuzzleGameStmt:   q.unlinkImageFromPuzzleGameStmt,
		unlockPuzzleGameStmt:            q.unlockPuzzleGameStmt,
		updatePuzzleGameStmt:            q.updatePuzzleGameStmt,
	}
}

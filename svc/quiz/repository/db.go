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
	if q.addNewPlayerStmt, err = db.PrepareContext(ctx, addNewPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query AddNewPlayer: %w", err)
	}
	if q.addNewQuizStmt, err = db.PrepareContext(ctx, addNewQuiz); err != nil {
		return nil, fmt.Errorf("error preparing query AddNewQuiz: %w", err)
	}
	if q.countCorrectAnswersStmt, err = db.PrepareContext(ctx, countCorrectAnswers); err != nil {
		return nil, fmt.Errorf("error preparing query CountCorrectAnswers: %w", err)
	}
	if q.countPlayersInQuizStmt, err = db.PrepareContext(ctx, countPlayersInQuiz); err != nil {
		return nil, fmt.Errorf("error preparing query CountPlayersInQuiz: %w", err)
	}
	if q.getQuizByChallengeIDStmt, err = db.PrepareContext(ctx, getQuizByChallengeID); err != nil {
		return nil, fmt.Errorf("error preparing query GetQuizByChallengeID: %w", err)
	}
	if q.getQuizByIDStmt, err = db.PrepareContext(ctx, getQuizByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetQuizByID: %w", err)
	}
	if q.getQuizWinnnersStmt, err = db.PrepareContext(ctx, getQuizWinnners); err != nil {
		return nil, fmt.Errorf("error preparing query GetQuizWinnners: %w", err)
	}
	if q.storeAnswerStmt, err = db.PrepareContext(ctx, storeAnswer); err != nil {
		return nil, fmt.Errorf("error preparing query StoreAnswer: %w", err)
	}
	if q.updatePlayerStatusStmt, err = db.PrepareContext(ctx, updatePlayerStatus); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePlayerStatus: %w", err)
	}
	if q.updateQuizStatusStmt, err = db.PrepareContext(ctx, updateQuizStatus); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateQuizStatus: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addNewPlayerStmt != nil {
		if cerr := q.addNewPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNewPlayerStmt: %w", cerr)
		}
	}
	if q.addNewQuizStmt != nil {
		if cerr := q.addNewQuizStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addNewQuizStmt: %w", cerr)
		}
	}
	if q.countCorrectAnswersStmt != nil {
		if cerr := q.countCorrectAnswersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countCorrectAnswersStmt: %w", cerr)
		}
	}
	if q.countPlayersInQuizStmt != nil {
		if cerr := q.countPlayersInQuizStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countPlayersInQuizStmt: %w", cerr)
		}
	}
	if q.getQuizByChallengeIDStmt != nil {
		if cerr := q.getQuizByChallengeIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getQuizByChallengeIDStmt: %w", cerr)
		}
	}
	if q.getQuizByIDStmt != nil {
		if cerr := q.getQuizByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getQuizByIDStmt: %w", cerr)
		}
	}
	if q.getQuizWinnnersStmt != nil {
		if cerr := q.getQuizWinnnersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getQuizWinnnersStmt: %w", cerr)
		}
	}
	if q.storeAnswerStmt != nil {
		if cerr := q.storeAnswerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing storeAnswerStmt: %w", cerr)
		}
	}
	if q.updatePlayerStatusStmt != nil {
		if cerr := q.updatePlayerStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePlayerStatusStmt: %w", cerr)
		}
	}
	if q.updateQuizStatusStmt != nil {
		if cerr := q.updateQuizStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateQuizStatusStmt: %w", cerr)
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
	db                       DBTX
	tx                       *sql.Tx
	addNewPlayerStmt         *sql.Stmt
	addNewQuizStmt           *sql.Stmt
	countCorrectAnswersStmt  *sql.Stmt
	countPlayersInQuizStmt   *sql.Stmt
	getQuizByChallengeIDStmt *sql.Stmt
	getQuizByIDStmt          *sql.Stmt
	getQuizWinnnersStmt      *sql.Stmt
	storeAnswerStmt          *sql.Stmt
	updatePlayerStatusStmt   *sql.Stmt
	updateQuizStatusStmt     *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                       tx,
		tx:                       tx,
		addNewPlayerStmt:         q.addNewPlayerStmt,
		addNewQuizStmt:           q.addNewQuizStmt,
		countCorrectAnswersStmt:  q.countCorrectAnswersStmt,
		countPlayersInQuizStmt:   q.countPlayersInQuizStmt,
		getQuizByChallengeIDStmt: q.getQuizByChallengeIDStmt,
		getQuizByIDStmt:          q.getQuizByIDStmt,
		getQuizWinnnersStmt:      q.getQuizWinnnersStmt,
		storeAnswerStmt:          q.storeAnswerStmt,
		updatePlayerStatusStmt:   q.updatePlayerStatusStmt,
		updateQuizStatusStmt:     q.updateQuizStatusStmt,
	}
}

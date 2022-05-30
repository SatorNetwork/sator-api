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
	if q.getTransactionsByStatusStmt, err = db.PrepareContext(ctx, getTransactionsByStatus); err != nil {
		return nil, fmt.Errorf("error preparing query GetTransactionsByStatus: %w", err)
	}
	if q.registerTransactionStmt, err = db.PrepareContext(ctx, registerTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query RegisterTransaction: %w", err)
	}
	if q.updateTransactionStmt, err = db.PrepareContext(ctx, updateTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTransaction: %w", err)
	}
	if q.updateTransactionStatusStmt, err = db.PrepareContext(ctx, updateTransactionStatus); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTransactionStatus: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.getTransactionsByStatusStmt != nil {
		if cerr := q.getTransactionsByStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTransactionsByStatusStmt: %w", cerr)
		}
	}
	if q.registerTransactionStmt != nil {
		if cerr := q.registerTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing registerTransactionStmt: %w", cerr)
		}
	}
	if q.updateTransactionStmt != nil {
		if cerr := q.updateTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTransactionStmt: %w", cerr)
		}
	}
	if q.updateTransactionStatusStmt != nil {
		if cerr := q.updateTransactionStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTransactionStatusStmt: %w", cerr)
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
	db                          DBTX
	tx                          *sql.Tx
	getTransactionsByStatusStmt *sql.Stmt
	registerTransactionStmt     *sql.Stmt
	updateTransactionStmt       *sql.Stmt
	updateTransactionStatusStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                          tx,
		tx:                          tx,
		getTransactionsByStatusStmt: q.getTransactionsByStatusStmt,
		registerTransactionStmt:     q.registerTransactionStmt,
		updateTransactionStmt:       q.updateTransactionStmt,
		updateTransactionStatusStmt: q.updateTransactionStatusStmt,
	}
}
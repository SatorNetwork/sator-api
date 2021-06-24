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
	if q.addTransactionStmt, err = db.PrepareContext(ctx, addTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query AddTransaction: %w", err)
	}
	if q.getTotalAmountStmt, err = db.PrepareContext(ctx, getTotalAmount); err != nil {
		return nil, fmt.Errorf("error preparing query GetTotalAmount: %w", err)
	}
	if q.getTransactionsByUserIDPaginatedStmt, err = db.PrepareContext(ctx, getTransactionsByUserIDPaginated); err != nil {
		return nil, fmt.Errorf("error preparing query GetTransactionsByUserIDPaginated: %w", err)
	}
	if q.withdrawStmt, err = db.PrepareContext(ctx, withdraw); err != nil {
		return nil, fmt.Errorf("error preparing query Withdraw: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addTransactionStmt != nil {
		if cerr := q.addTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addTransactionStmt: %w", cerr)
		}
	}
	if q.getTotalAmountStmt != nil {
		if cerr := q.getTotalAmountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTotalAmountStmt: %w", cerr)
		}
	}
	if q.getTransactionsByUserIDPaginatedStmt != nil {
		if cerr := q.getTransactionsByUserIDPaginatedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTransactionsByUserIDPaginatedStmt: %w", cerr)
		}
	}
	if q.withdrawStmt != nil {
		if cerr := q.withdrawStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing withdrawStmt: %w", cerr)
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
	db                                   DBTX
	tx                                   *sql.Tx
	addTransactionStmt                   *sql.Stmt
	getTotalAmountStmt                   *sql.Stmt
	getTransactionsByUserIDPaginatedStmt *sql.Stmt
	withdrawStmt                         *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                   tx,
		tx:                                   tx,
		addTransactionStmt:                   q.addTransactionStmt,
		getTotalAmountStmt:                   q.getTotalAmountStmt,
		getTransactionsByUserIDPaginatedStmt: q.getTransactionsByUserIDPaginatedStmt,
		withdrawStmt:                         q.withdrawStmt,
	}
}

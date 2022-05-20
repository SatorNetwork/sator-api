// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

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
	if q.getAmountAvailableToWithdrawStmt, err = db.PrepareContext(ctx, getAmountAvailableToWithdraw); err != nil {
		return nil, fmt.Errorf("error preparing query GetAmountAvailableToWithdraw: %w", err)
	}
	if q.getFailedTransactionsStmt, err = db.PrepareContext(ctx, getFailedTransactions); err != nil {
		return nil, fmt.Errorf("error preparing query GetFailedTransactions: %w", err)
	}
	if q.getRequestedTransactionsStmt, err = db.PrepareContext(ctx, getRequestedTransactions); err != nil {
		return nil, fmt.Errorf("error preparing query GetRequestedTransactions: %w", err)
	}
	if q.getScannedQRCodeByUserIDStmt, err = db.PrepareContext(ctx, getScannedQRCodeByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetScannedQRCodeByUserID: %w", err)
	}
	if q.getTotalAmountStmt, err = db.PrepareContext(ctx, getTotalAmount); err != nil {
		return nil, fmt.Errorf("error preparing query GetTotalAmount: %w", err)
	}
	if q.getTransactionsByUserIDPaginatedStmt, err = db.PrepareContext(ctx, getTransactionsByUserIDPaginated); err != nil {
		return nil, fmt.Errorf("error preparing query GetTransactionsByUserIDPaginated: %w", err)
	}
	if q.requestTransactionsByUserIDStmt, err = db.PrepareContext(ctx, requestTransactionsByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query RequestTransactionsByUserID: %w", err)
	}
	if q.updateTransactionStatusByTxHashStmt, err = db.PrepareContext(ctx, updateTransactionStatusByTxHash); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTransactionStatusByTxHash: %w", err)
	}
	if q.updateTransactionTxHashStmt, err = db.PrepareContext(ctx, updateTransactionTxHash); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTransactionTxHash: %w", err)
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
	if q.getAmountAvailableToWithdrawStmt != nil {
		if cerr := q.getAmountAvailableToWithdrawStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAmountAvailableToWithdrawStmt: %w", cerr)
		}
	}
	if q.getFailedTransactionsStmt != nil {
		if cerr := q.getFailedTransactionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFailedTransactionsStmt: %w", cerr)
		}
	}
	if q.getRequestedTransactionsStmt != nil {
		if cerr := q.getRequestedTransactionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRequestedTransactionsStmt: %w", cerr)
		}
	}
	if q.getScannedQRCodeByUserIDStmt != nil {
		if cerr := q.getScannedQRCodeByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getScannedQRCodeByUserIDStmt: %w", cerr)
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
	if q.requestTransactionsByUserIDStmt != nil {
		if cerr := q.requestTransactionsByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing requestTransactionsByUserIDStmt: %w", cerr)
		}
	}
	if q.updateTransactionStatusByTxHashStmt != nil {
		if cerr := q.updateTransactionStatusByTxHashStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTransactionStatusByTxHashStmt: %w", cerr)
		}
	}
	if q.updateTransactionTxHashStmt != nil {
		if cerr := q.updateTransactionTxHashStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTransactionTxHashStmt: %w", cerr)
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
	getAmountAvailableToWithdrawStmt     *sql.Stmt
	getFailedTransactionsStmt            *sql.Stmt
	getRequestedTransactionsStmt         *sql.Stmt
	getScannedQRCodeByUserIDStmt         *sql.Stmt
	getTotalAmountStmt                   *sql.Stmt
	getTransactionsByUserIDPaginatedStmt *sql.Stmt
	requestTransactionsByUserIDStmt      *sql.Stmt
	updateTransactionStatusByTxHashStmt  *sql.Stmt
	updateTransactionTxHashStmt          *sql.Stmt
	withdrawStmt                         *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                   tx,
		tx:                                   tx,
		addTransactionStmt:                   q.addTransactionStmt,
		getAmountAvailableToWithdrawStmt:     q.getAmountAvailableToWithdrawStmt,
		getFailedTransactionsStmt:            q.getFailedTransactionsStmt,
		getRequestedTransactionsStmt:         q.getRequestedTransactionsStmt,
		getScannedQRCodeByUserIDStmt:         q.getScannedQRCodeByUserIDStmt,
		getTotalAmountStmt:                   q.getTotalAmountStmt,
		getTransactionsByUserIDPaginatedStmt: q.getTransactionsByUserIDPaginatedStmt,
		requestTransactionsByUserIDStmt:      q.requestTransactionsByUserIDStmt,
		updateTransactionStatusByTxHashStmt:  q.updateTransactionStatusByTxHashStmt,
		updateTransactionTxHashStmt:          q.updateTransactionTxHashStmt,
		withdrawStmt:                         q.withdrawStmt,
	}
}

package db

import (
	"context"
	"database/sql"
)

type (
	// TransactionFunc type
	TransactionFunc func(txFunc func(DBTX) error) (err error)

	// DBTX ...
	// Database transaction interface
	DBTX interface {
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		PrepareContext(context.Context, string) (*sql.Stmt, error)
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}
)

// Transaction is a wrapper function which helps to avoid the use of sql.DB instance directly
func Transaction(db *sql.DB) TransactionFunc {
	return func(txFunc func(DBTX) error) (err error) {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p) // re-throw panic after Rollback
			} else if err != nil {
				tx.Rollback() // err is non-nil; don't change it
			} else {
				err = tx.Commit() // err is nil; if Commit returns error update err
			}
		}()

		return txFunc(tx)
	}
}

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
	if q.addSettingStmt, err = db.PrepareContext(ctx, addSetting); err != nil {
		return nil, fmt.Errorf("error preparing query AddSetting: %w", err)
	}
	if q.deleteSettingStmt, err = db.PrepareContext(ctx, deleteSetting); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSetting: %w", err)
	}
	if q.getSettingByKeyStmt, err = db.PrepareContext(ctx, getSettingByKey); err != nil {
		return nil, fmt.Errorf("error preparing query GetSettingByKey: %w", err)
	}
	if q.getSettingsStmt, err = db.PrepareContext(ctx, getSettings); err != nil {
		return nil, fmt.Errorf("error preparing query GetSettings: %w", err)
	}
	if q.updateSettingStmt, err = db.PrepareContext(ctx, updateSetting); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSetting: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addSettingStmt != nil {
		if cerr := q.addSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addSettingStmt: %w", cerr)
		}
	}
	if q.deleteSettingStmt != nil {
		if cerr := q.deleteSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSettingStmt: %w", cerr)
		}
	}
	if q.getSettingByKeyStmt != nil {
		if cerr := q.getSettingByKeyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSettingByKeyStmt: %w", cerr)
		}
	}
	if q.getSettingsStmt != nil {
		if cerr := q.getSettingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSettingsStmt: %w", cerr)
		}
	}
	if q.updateSettingStmt != nil {
		if cerr := q.updateSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSettingStmt: %w", cerr)
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
	db                  DBTX
	tx                  *sql.Tx
	addSettingStmt      *sql.Stmt
	deleteSettingStmt   *sql.Stmt
	getSettingByKeyStmt *sql.Stmt
	getSettingsStmt     *sql.Stmt
	updateSettingStmt   *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                  tx,
		tx:                  tx,
		addSettingStmt:      q.addSettingStmt,
		deleteSettingStmt:   q.deleteSettingStmt,
		getSettingByKeyStmt: q.getSettingByKeyStmt,
		getSettingsStmt:     q.getSettingsStmt,
		updateSettingStmt:   q.updateSettingStmt,
	}
}

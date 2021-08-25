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
	if q.addReferralStmt, err = db.PrepareContext(ctx, addReferral); err != nil {
		return nil, fmt.Errorf("error preparing query AddReferral: %w", err)
	}
	if q.addReferralCodeDataStmt, err = db.PrepareContext(ctx, addReferralCodeData); err != nil {
		return nil, fmt.Errorf("error preparing query AddReferralCodeData: %w", err)
	}
	if q.deleteReferralCodeDataByIDStmt, err = db.PrepareContext(ctx, deleteReferralCodeDataByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteReferralCodeDataByID: %w", err)
	}
	if q.getReferralCodeDataByUserIDStmt, err = db.PrepareContext(ctx, getReferralCodeDataByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetReferralCodeDataByUserID: %w", err)
	}
	if q.getReferralCodesDataListStmt, err = db.PrepareContext(ctx, getReferralCodesDataList); err != nil {
		return nil, fmt.Errorf("error preparing query GetReferralCodesDataList: %w", err)
	}
	if q.getReferralsWithPaginationByUserIDStmt, err = db.PrepareContext(ctx, getReferralsWithPaginationByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetReferralsWithPaginationByUserID: %w", err)
	}
	if q.updateReferralCodeDataStmt, err = db.PrepareContext(ctx, updateReferralCodeData); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateReferralCodeData: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addReferralStmt != nil {
		if cerr := q.addReferralStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addReferralStmt: %w", cerr)
		}
	}
	if q.addReferralCodeDataStmt != nil {
		if cerr := q.addReferralCodeDataStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addReferralCodeDataStmt: %w", cerr)
		}
	}
	if q.deleteReferralCodeDataByIDStmt != nil {
		if cerr := q.deleteReferralCodeDataByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteReferralCodeDataByIDStmt: %w", cerr)
		}
	}
	if q.getReferralCodeDataByUserIDStmt != nil {
		if cerr := q.getReferralCodeDataByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getReferralCodeDataByUserIDStmt: %w", cerr)
		}
	}
	if q.getReferralCodesDataListStmt != nil {
		if cerr := q.getReferralCodesDataListStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getReferralCodesDataListStmt: %w", cerr)
		}
	}
	if q.getReferralsWithPaginationByUserIDStmt != nil {
		if cerr := q.getReferralsWithPaginationByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getReferralsWithPaginationByUserIDStmt: %w", cerr)
		}
	}
	if q.updateReferralCodeDataStmt != nil {
		if cerr := q.updateReferralCodeDataStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateReferralCodeDataStmt: %w", cerr)
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
	addReferralStmt                        *sql.Stmt
	addReferralCodeDataStmt                *sql.Stmt
	deleteReferralCodeDataByIDStmt         *sql.Stmt
	getReferralCodeDataByUserIDStmt        *sql.Stmt
	getReferralCodesDataListStmt           *sql.Stmt
	getReferralsWithPaginationByUserIDStmt *sql.Stmt
	updateReferralCodeDataStmt             *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                     tx,
		tx:                                     tx,
		addReferralStmt:                        q.addReferralStmt,
		addReferralCodeDataStmt:                q.addReferralCodeDataStmt,
		deleteReferralCodeDataByIDStmt:         q.deleteReferralCodeDataByIDStmt,
		getReferralCodeDataByUserIDStmt:        q.getReferralCodeDataByUserIDStmt,
		getReferralCodesDataListStmt:           q.getReferralCodesDataListStmt,
		getReferralsWithPaginationByUserIDStmt: q.getReferralsWithPaginationByUserIDStmt,
		updateReferralCodeDataStmt:             q.updateReferralCodeDataStmt,
	}
}

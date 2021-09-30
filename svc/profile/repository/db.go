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
	if q.createProfileStmt, err = db.PrepareContext(ctx, createProfile); err != nil {
		return nil, fmt.Errorf("error preparing query CreateProfile: %w", err)
	}
	if q.getProfileByUserIDStmt, err = db.PrepareContext(ctx, getProfileByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfileByUserID: %w", err)
	}
	if q.updateAvatarStmt, err = db.PrepareContext(ctx, updateAvatar); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAvatar: %w", err)
	}
	if q.updateProfileByIDStmt, err = db.PrepareContext(ctx, updateProfileByID); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfileByID: %w", err)
	}
	if q.updateProfileByUserIDStmt, err = db.PrepareContext(ctx, updateProfileByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfileByUserID: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createProfileStmt != nil {
		if cerr := q.createProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createProfileStmt: %w", cerr)
		}
	}
	if q.getProfileByUserIDStmt != nil {
		if cerr := q.getProfileByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileByUserIDStmt: %w", cerr)
		}
	}
	if q.updateAvatarStmt != nil {
		if cerr := q.updateAvatarStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAvatarStmt: %w", cerr)
		}
	}
	if q.updateProfileByIDStmt != nil {
		if cerr := q.updateProfileByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileByIDStmt: %w", cerr)
		}
	}
	if q.updateProfileByUserIDStmt != nil {
		if cerr := q.updateProfileByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileByUserIDStmt: %w", cerr)
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
	db                        DBTX
	tx                        *sql.Tx
	createProfileStmt         *sql.Stmt
	getProfileByUserIDStmt    *sql.Stmt
	updateAvatarStmt          *sql.Stmt
	updateProfileByIDStmt     *sql.Stmt
	updateProfileByUserIDStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                        tx,
		tx:                        tx,
		createProfileStmt:         q.createProfileStmt,
		getProfileByUserIDStmt:    q.getProfileByUserIDStmt,
		updateAvatarStmt:          q.updateAvatarStmt,
		updateProfileByIDStmt:     q.updateProfileByIDStmt,
		updateProfileByUserIDStmt: q.updateProfileByUserIDStmt,
	}
}

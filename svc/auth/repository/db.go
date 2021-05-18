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
	if q.createPasswordResetStmt, err = db.PrepareContext(ctx, createPasswordReset); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePasswordReset: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.createUserVerificationStmt, err = db.PrepareContext(ctx, createUserVerification); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUserVerification: %w", err)
	}
	if q.deletePasswordResetsByEmailStmt, err = db.PrepareContext(ctx, deletePasswordResetsByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query DeletePasswordResetsByEmail: %w", err)
	}
	if q.deletePasswordResetsByUserIDStmt, err = db.PrepareContext(ctx, deletePasswordResetsByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query DeletePasswordResetsByUserID: %w", err)
	}
	if q.deleteUserByIDStmt, err = db.PrepareContext(ctx, deleteUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserByID: %w", err)
	}
	if q.deleteUserVerificationsByEmailStmt, err = db.PrepareContext(ctx, deleteUserVerificationsByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserVerificationsByEmail: %w", err)
	}
	if q.deleteUserVerificationsByUserIDStmt, err = db.PrepareContext(ctx, deleteUserVerificationsByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserVerificationsByUserID: %w", err)
	}
	if q.getPasswordResetByEmailStmt, err = db.PrepareContext(ctx, getPasswordResetByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetPasswordResetByEmail: %w", err)
	}
	if q.getUserByEmailStmt, err = db.PrepareContext(ctx, getUserByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByEmail: %w", err)
	}
	if q.getUserByIDStmt, err = db.PrepareContext(ctx, getUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByID: %w", err)
	}
	if q.getUserVerificationByUserIDStmt, err = db.PrepareContext(ctx, getUserVerificationByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserVerificationByUserID: %w", err)
	}
	if q.getUsersListDescStmt, err = db.PrepareContext(ctx, getUsersListDesc); err != nil {
		return nil, fmt.Errorf("error preparing query GetUsersListDesc: %w", err)
	}
	if q.updateUserEmailStmt, err = db.PrepareContext(ctx, updateUserEmail); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserEmail: %w", err)
	}
	if q.updateUserPasswordStmt, err = db.PrepareContext(ctx, updateUserPassword); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserPassword: %w", err)
	}
	if q.updateUserStatusStmt, err = db.PrepareContext(ctx, updateUserStatus); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserStatus: %w", err)
	}
	if q.updateUserVerifiedAtStmt, err = db.PrepareContext(ctx, updateUserVerifiedAt); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserVerifiedAt: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createPasswordResetStmt != nil {
		if cerr := q.createPasswordResetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPasswordResetStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.createUserVerificationStmt != nil {
		if cerr := q.createUserVerificationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserVerificationStmt: %w", cerr)
		}
	}
	if q.deletePasswordResetsByEmailStmt != nil {
		if cerr := q.deletePasswordResetsByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deletePasswordResetsByEmailStmt: %w", cerr)
		}
	}
	if q.deletePasswordResetsByUserIDStmt != nil {
		if cerr := q.deletePasswordResetsByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deletePasswordResetsByUserIDStmt: %w", cerr)
		}
	}
	if q.deleteUserByIDStmt != nil {
		if cerr := q.deleteUserByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserByIDStmt: %w", cerr)
		}
	}
	if q.deleteUserVerificationsByEmailStmt != nil {
		if cerr := q.deleteUserVerificationsByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserVerificationsByEmailStmt: %w", cerr)
		}
	}
	if q.deleteUserVerificationsByUserIDStmt != nil {
		if cerr := q.deleteUserVerificationsByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserVerificationsByUserIDStmt: %w", cerr)
		}
	}
	if q.getPasswordResetByEmailStmt != nil {
		if cerr := q.getPasswordResetByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPasswordResetByEmailStmt: %w", cerr)
		}
	}
	if q.getUserByEmailStmt != nil {
		if cerr := q.getUserByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByEmailStmt: %w", cerr)
		}
	}
	if q.getUserByIDStmt != nil {
		if cerr := q.getUserByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByIDStmt: %w", cerr)
		}
	}
	if q.getUserVerificationByUserIDStmt != nil {
		if cerr := q.getUserVerificationByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserVerificationByUserIDStmt: %w", cerr)
		}
	}
	if q.getUsersListDescStmt != nil {
		if cerr := q.getUsersListDescStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUsersListDescStmt: %w", cerr)
		}
	}
	if q.updateUserEmailStmt != nil {
		if cerr := q.updateUserEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserEmailStmt: %w", cerr)
		}
	}
	if q.updateUserPasswordStmt != nil {
		if cerr := q.updateUserPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserPasswordStmt: %w", cerr)
		}
	}
	if q.updateUserStatusStmt != nil {
		if cerr := q.updateUserStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserStatusStmt: %w", cerr)
		}
	}
	if q.updateUserVerifiedAtStmt != nil {
		if cerr := q.updateUserVerifiedAtStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserVerifiedAtStmt: %w", cerr)
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
	db                                  DBTX
	tx                                  *sql.Tx
	createPasswordResetStmt             *sql.Stmt
	createUserStmt                      *sql.Stmt
	createUserVerificationStmt          *sql.Stmt
	deletePasswordResetsByEmailStmt     *sql.Stmt
	deletePasswordResetsByUserIDStmt    *sql.Stmt
	deleteUserByIDStmt                  *sql.Stmt
	deleteUserVerificationsByEmailStmt  *sql.Stmt
	deleteUserVerificationsByUserIDStmt *sql.Stmt
	getPasswordResetByEmailStmt         *sql.Stmt
	getUserByEmailStmt                  *sql.Stmt
	getUserByIDStmt                     *sql.Stmt
	getUserVerificationByUserIDStmt     *sql.Stmt
	getUsersListDescStmt                *sql.Stmt
	updateUserEmailStmt                 *sql.Stmt
	updateUserPasswordStmt              *sql.Stmt
	updateUserStatusStmt                *sql.Stmt
	updateUserVerifiedAtStmt            *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                  tx,
		tx:                                  tx,
		createPasswordResetStmt:             q.createPasswordResetStmt,
		createUserStmt:                      q.createUserStmt,
		createUserVerificationStmt:          q.createUserVerificationStmt,
		deletePasswordResetsByEmailStmt:     q.deletePasswordResetsByEmailStmt,
		deletePasswordResetsByUserIDStmt:    q.deletePasswordResetsByUserIDStmt,
		deleteUserByIDStmt:                  q.deleteUserByIDStmt,
		deleteUserVerificationsByEmailStmt:  q.deleteUserVerificationsByEmailStmt,
		deleteUserVerificationsByUserIDStmt: q.deleteUserVerificationsByUserIDStmt,
		getPasswordResetByEmailStmt:         q.getPasswordResetByEmailStmt,
		getUserByEmailStmt:                  q.getUserByEmailStmt,
		getUserByIDStmt:                     q.getUserByIDStmt,
		getUserVerificationByUserIDStmt:     q.getUserVerificationByUserIDStmt,
		getUsersListDescStmt:                q.getUsersListDescStmt,
		updateUserEmailStmt:                 q.updateUserEmailStmt,
		updateUserPasswordStmt:              q.updateUserPasswordStmt,
		updateUserStatusStmt:                q.updateUserStatusStmt,
		updateUserVerifiedAtStmt:            q.updateUserVerifiedAtStmt,
	}
}

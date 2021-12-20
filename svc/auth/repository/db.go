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
	if q.addToBlacklistStmt, err = db.PrepareContext(ctx, addToBlacklist); err != nil {
		return nil, fmt.Errorf("error preparing query AddToBlacklist: %w", err)
	}
	if q.addToWhitelistStmt, err = db.PrepareContext(ctx, addToWhitelist); err != nil {
		return nil, fmt.Errorf("error preparing query AddToWhitelist: %w", err)
	}
	if q.blockUsersOnTheSameDeviceStmt, err = db.PrepareContext(ctx, blockUsersOnTheSameDevice); err != nil {
		return nil, fmt.Errorf("error preparing query BlockUsersOnTheSameDevice: %w", err)
	}
	if q.blockUsersWithDuplicateEmailStmt, err = db.PrepareContext(ctx, blockUsersWithDuplicateEmail); err != nil {
		return nil, fmt.Errorf("error preparing query BlockUsersWithDuplicateEmail: %w", err)
	}
	if q.countAllUsersStmt, err = db.PrepareContext(ctx, countAllUsers); err != nil {
		return nil, fmt.Errorf("error preparing query CountAllUsers: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.createUserVerificationStmt, err = db.PrepareContext(ctx, createUserVerification); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUserVerification: %w", err)
	}
	if q.deleteFromBlacklistStmt, err = db.PrepareContext(ctx, deleteFromBlacklist); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteFromBlacklist: %w", err)
	}
	if q.deleteFromWhitelistStmt, err = db.PrepareContext(ctx, deleteFromWhitelist); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteFromWhitelist: %w", err)
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
	if q.destroyUserStmt, err = db.PrepareContext(ctx, destroyUser); err != nil {
		return nil, fmt.Errorf("error preparing query DestroyUser: %w", err)
	}
	if q.doesUserHaveMoreThanOneAccountStmt, err = db.PrepareContext(ctx, doesUserHaveMoreThanOneAccount); err != nil {
		return nil, fmt.Errorf("error preparing query DoesUserHaveMoreThanOneAccount: %w", err)
	}
	if q.getBlacklistStmt, err = db.PrepareContext(ctx, getBlacklist); err != nil {
		return nil, fmt.Errorf("error preparing query GetBlacklist: %w", err)
	}
	if q.getBlacklistByRestrictedValueStmt, err = db.PrepareContext(ctx, getBlacklistByRestrictedValue); err != nil {
		return nil, fmt.Errorf("error preparing query GetBlacklistByRestrictedValue: %w", err)
	}
	if q.getKYCStatusStmt, err = db.PrepareContext(ctx, getKYCStatus); err != nil {
		return nil, fmt.Errorf("error preparing query GetKYCStatus: %w", err)
	}
	if q.getNotSanitizedUsersListDescStmt, err = db.PrepareContext(ctx, getNotSanitizedUsersListDesc); err != nil {
		return nil, fmt.Errorf("error preparing query GetNotSanitizedUsersListDesc: %w", err)
	}
	if q.getUserByEmailStmt, err = db.PrepareContext(ctx, getUserByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByEmail: %w", err)
	}
	if q.getUserByIDStmt, err = db.PrepareContext(ctx, getUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByID: %w", err)
	}
	if q.getUserBySanitizedEmailStmt, err = db.PrepareContext(ctx, getUserBySanitizedEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserBySanitizedEmail: %w", err)
	}
	if q.getUserByUsernameStmt, err = db.PrepareContext(ctx, getUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByUsername: %w", err)
	}
	if q.getUserIDsOnTheSameDeviceStmt, err = db.PrepareContext(ctx, getUserIDsOnTheSameDevice); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserIDsOnTheSameDevice: %w", err)
	}
	if q.getUserVerificationByEmailStmt, err = db.PrepareContext(ctx, getUserVerificationByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserVerificationByEmail: %w", err)
	}
	if q.getUserVerificationByUserIDStmt, err = db.PrepareContext(ctx, getUserVerificationByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserVerificationByUserID: %w", err)
	}
	if q.getUsernameByIDStmt, err = db.PrepareContext(ctx, getUsernameByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUsernameByID: %w", err)
	}
	if q.getUsersListDescStmt, err = db.PrepareContext(ctx, getUsersListDesc); err != nil {
		return nil, fmt.Errorf("error preparing query GetUsersListDesc: %w", err)
	}
	if q.getVerifiedUsersListDescStmt, err = db.PrepareContext(ctx, getVerifiedUsersListDesc); err != nil {
		return nil, fmt.Errorf("error preparing query GetVerifiedUsersListDesc: %w", err)
	}
	if q.getWhitelistStmt, err = db.PrepareContext(ctx, getWhitelist); err != nil {
		return nil, fmt.Errorf("error preparing query GetWhitelist: %w", err)
	}
	if q.getWhitelistByAllowedValueStmt, err = db.PrepareContext(ctx, getWhitelistByAllowedValue); err != nil {
		return nil, fmt.Errorf("error preparing query GetWhitelistByAllowedValue: %w", err)
	}
	if q.isEmailBlacklistedStmt, err = db.PrepareContext(ctx, isEmailBlacklisted); err != nil {
		return nil, fmt.Errorf("error preparing query IsEmailBlacklisted: %w", err)
	}
	if q.isEmailWhitelistedStmt, err = db.PrepareContext(ctx, isEmailWhitelisted); err != nil {
		return nil, fmt.Errorf("error preparing query IsEmailWhitelisted: %w", err)
	}
	if q.isUserDisabledStmt, err = db.PrepareContext(ctx, isUserDisabled); err != nil {
		return nil, fmt.Errorf("error preparing query IsUserDisabled: %w", err)
	}
	if q.linkDeviceToUserStmt, err = db.PrepareContext(ctx, linkDeviceToUser); err != nil {
		return nil, fmt.Errorf("error preparing query LinkDeviceToUser: %w", err)
	}
	if q.updateKYCStatusStmt, err = db.PrepareContext(ctx, updateKYCStatus); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateKYCStatus: %w", err)
	}
	if q.updatePublicKeyStmt, err = db.PrepareContext(ctx, updatePublicKey); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePublicKey: %w", err)
	}
	if q.updateUserEmailStmt, err = db.PrepareContext(ctx, updateUserEmail); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserEmail: %w", err)
	}
	if q.updateUserPasswordStmt, err = db.PrepareContext(ctx, updateUserPassword); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserPassword: %w", err)
	}
	if q.updateUserSanitizedEmailStmt, err = db.PrepareContext(ctx, updateUserSanitizedEmail); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserSanitizedEmail: %w", err)
	}
	if q.updateUserStatusStmt, err = db.PrepareContext(ctx, updateUserStatus); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserStatus: %w", err)
	}
	if q.updateUserVerifiedAtStmt, err = db.PrepareContext(ctx, updateUserVerifiedAt); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserVerifiedAt: %w", err)
	}
	if q.updateUsernameStmt, err = db.PrepareContext(ctx, updateUsername); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUsername: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addToBlacklistStmt != nil {
		if cerr := q.addToBlacklistStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addToBlacklistStmt: %w", cerr)
		}
	}
	if q.addToWhitelistStmt != nil {
		if cerr := q.addToWhitelistStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addToWhitelistStmt: %w", cerr)
		}
	}
	if q.blockUsersOnTheSameDeviceStmt != nil {
		if cerr := q.blockUsersOnTheSameDeviceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing blockUsersOnTheSameDeviceStmt: %w", cerr)
		}
	}
	if q.blockUsersWithDuplicateEmailStmt != nil {
		if cerr := q.blockUsersWithDuplicateEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing blockUsersWithDuplicateEmailStmt: %w", cerr)
		}
	}
	if q.countAllUsersStmt != nil {
		if cerr := q.countAllUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countAllUsersStmt: %w", cerr)
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
	if q.deleteFromBlacklistStmt != nil {
		if cerr := q.deleteFromBlacklistStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteFromBlacklistStmt: %w", cerr)
		}
	}
	if q.deleteFromWhitelistStmt != nil {
		if cerr := q.deleteFromWhitelistStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteFromWhitelistStmt: %w", cerr)
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
	if q.destroyUserStmt != nil {
		if cerr := q.destroyUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing destroyUserStmt: %w", cerr)
		}
	}
	if q.doesUserHaveMoreThanOneAccountStmt != nil {
		if cerr := q.doesUserHaveMoreThanOneAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing doesUserHaveMoreThanOneAccountStmt: %w", cerr)
		}
	}
	if q.getBlacklistStmt != nil {
		if cerr := q.getBlacklistStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBlacklistStmt: %w", cerr)
		}
	}
	if q.getBlacklistByRestrictedValueStmt != nil {
		if cerr := q.getBlacklistByRestrictedValueStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBlacklistByRestrictedValueStmt: %w", cerr)
		}
	}
	if q.getKYCStatusStmt != nil {
		if cerr := q.getKYCStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getKYCStatusStmt: %w", cerr)
		}
	}
	if q.getNotSanitizedUsersListDescStmt != nil {
		if cerr := q.getNotSanitizedUsersListDescStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNotSanitizedUsersListDescStmt: %w", cerr)
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
	if q.getUserBySanitizedEmailStmt != nil {
		if cerr := q.getUserBySanitizedEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserBySanitizedEmailStmt: %w", cerr)
		}
	}
	if q.getUserByUsernameStmt != nil {
		if cerr := q.getUserByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByUsernameStmt: %w", cerr)
		}
	}
	if q.getUserIDsOnTheSameDeviceStmt != nil {
		if cerr := q.getUserIDsOnTheSameDeviceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserIDsOnTheSameDeviceStmt: %w", cerr)
		}
	}
	if q.getUserVerificationByEmailStmt != nil {
		if cerr := q.getUserVerificationByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserVerificationByEmailStmt: %w", cerr)
		}
	}
	if q.getUserVerificationByUserIDStmt != nil {
		if cerr := q.getUserVerificationByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserVerificationByUserIDStmt: %w", cerr)
		}
	}
	if q.getUsernameByIDStmt != nil {
		if cerr := q.getUsernameByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUsernameByIDStmt: %w", cerr)
		}
	}
	if q.getUsersListDescStmt != nil {
		if cerr := q.getUsersListDescStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUsersListDescStmt: %w", cerr)
		}
	}
	if q.getVerifiedUsersListDescStmt != nil {
		if cerr := q.getVerifiedUsersListDescStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getVerifiedUsersListDescStmt: %w", cerr)
		}
	}
	if q.getWhitelistStmt != nil {
		if cerr := q.getWhitelistStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWhitelistStmt: %w", cerr)
		}
	}
	if q.getWhitelistByAllowedValueStmt != nil {
		if cerr := q.getWhitelistByAllowedValueStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWhitelistByAllowedValueStmt: %w", cerr)
		}
	}
	if q.isEmailBlacklistedStmt != nil {
		if cerr := q.isEmailBlacklistedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing isEmailBlacklistedStmt: %w", cerr)
		}
	}
	if q.isEmailWhitelistedStmt != nil {
		if cerr := q.isEmailWhitelistedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing isEmailWhitelistedStmt: %w", cerr)
		}
	}
	if q.isUserDisabledStmt != nil {
		if cerr := q.isUserDisabledStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing isUserDisabledStmt: %w", cerr)
		}
	}
	if q.linkDeviceToUserStmt != nil {
		if cerr := q.linkDeviceToUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing linkDeviceToUserStmt: %w", cerr)
		}
	}
	if q.updateKYCStatusStmt != nil {
		if cerr := q.updateKYCStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateKYCStatusStmt: %w", cerr)
		}
	}
	if q.updatePublicKeyStmt != nil {
		if cerr := q.updatePublicKeyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePublicKeyStmt: %w", cerr)
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
	if q.updateUserSanitizedEmailStmt != nil {
		if cerr := q.updateUserSanitizedEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserSanitizedEmailStmt: %w", cerr)
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
	if q.updateUsernameStmt != nil {
		if cerr := q.updateUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUsernameStmt: %w", cerr)
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
	addToBlacklistStmt                  *sql.Stmt
	addToWhitelistStmt                  *sql.Stmt
	blockUsersOnTheSameDeviceStmt       *sql.Stmt
	blockUsersWithDuplicateEmailStmt    *sql.Stmt
	countAllUsersStmt                   *sql.Stmt
	createUserStmt                      *sql.Stmt
	createUserVerificationStmt          *sql.Stmt
	deleteFromBlacklistStmt             *sql.Stmt
	deleteFromWhitelistStmt             *sql.Stmt
	deleteUserByIDStmt                  *sql.Stmt
	deleteUserVerificationsByEmailStmt  *sql.Stmt
	deleteUserVerificationsByUserIDStmt *sql.Stmt
	destroyUserStmt                     *sql.Stmt
	doesUserHaveMoreThanOneAccountStmt  *sql.Stmt
	getBlacklistStmt                    *sql.Stmt
	getBlacklistByRestrictedValueStmt   *sql.Stmt
	getKYCStatusStmt                    *sql.Stmt
	getNotSanitizedUsersListDescStmt    *sql.Stmt
	getUserByEmailStmt                  *sql.Stmt
	getUserByIDStmt                     *sql.Stmt
	getUserBySanitizedEmailStmt         *sql.Stmt
	getUserByUsernameStmt               *sql.Stmt
	getUserIDsOnTheSameDeviceStmt       *sql.Stmt
	getUserVerificationByEmailStmt      *sql.Stmt
	getUserVerificationByUserIDStmt     *sql.Stmt
	getUsernameByIDStmt                 *sql.Stmt
	getUsersListDescStmt                *sql.Stmt
	getVerifiedUsersListDescStmt        *sql.Stmt
	getWhitelistStmt                    *sql.Stmt
	getWhitelistByAllowedValueStmt      *sql.Stmt
	isEmailBlacklistedStmt              *sql.Stmt
	isEmailWhitelistedStmt              *sql.Stmt
	isUserDisabledStmt                  *sql.Stmt
	linkDeviceToUserStmt                *sql.Stmt
	updateKYCStatusStmt                 *sql.Stmt
	updatePublicKeyStmt                 *sql.Stmt
	updateUserEmailStmt                 *sql.Stmt
	updateUserPasswordStmt              *sql.Stmt
	updateUserSanitizedEmailStmt        *sql.Stmt
	updateUserStatusStmt                *sql.Stmt
	updateUserVerifiedAtStmt            *sql.Stmt
	updateUsernameStmt                  *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                  tx,
		tx:                                  tx,
		addToBlacklistStmt:                  q.addToBlacklistStmt,
		addToWhitelistStmt:                  q.addToWhitelistStmt,
		blockUsersOnTheSameDeviceStmt:       q.blockUsersOnTheSameDeviceStmt,
		blockUsersWithDuplicateEmailStmt:    q.blockUsersWithDuplicateEmailStmt,
		countAllUsersStmt:                   q.countAllUsersStmt,
		createUserStmt:                      q.createUserStmt,
		createUserVerificationStmt:          q.createUserVerificationStmt,
		deleteFromBlacklistStmt:             q.deleteFromBlacklistStmt,
		deleteFromWhitelistStmt:             q.deleteFromWhitelistStmt,
		deleteUserByIDStmt:                  q.deleteUserByIDStmt,
		deleteUserVerificationsByEmailStmt:  q.deleteUserVerificationsByEmailStmt,
		deleteUserVerificationsByUserIDStmt: q.deleteUserVerificationsByUserIDStmt,
		destroyUserStmt:                     q.destroyUserStmt,
		doesUserHaveMoreThanOneAccountStmt:  q.doesUserHaveMoreThanOneAccountStmt,
		getBlacklistStmt:                    q.getBlacklistStmt,
		getBlacklistByRestrictedValueStmt:   q.getBlacklistByRestrictedValueStmt,
		getKYCStatusStmt:                    q.getKYCStatusStmt,
		getNotSanitizedUsersListDescStmt:    q.getNotSanitizedUsersListDescStmt,
		getUserByEmailStmt:                  q.getUserByEmailStmt,
		getUserByIDStmt:                     q.getUserByIDStmt,
		getUserBySanitizedEmailStmt:         q.getUserBySanitizedEmailStmt,
		getUserByUsernameStmt:               q.getUserByUsernameStmt,
		getUserIDsOnTheSameDeviceStmt:       q.getUserIDsOnTheSameDeviceStmt,
		getUserVerificationByEmailStmt:      q.getUserVerificationByEmailStmt,
		getUserVerificationByUserIDStmt:     q.getUserVerificationByUserIDStmt,
		getUsernameByIDStmt:                 q.getUsernameByIDStmt,
		getUsersListDescStmt:                q.getUsersListDescStmt,
		getVerifiedUsersListDescStmt:        q.getVerifiedUsersListDescStmt,
		getWhitelistStmt:                    q.getWhitelistStmt,
		getWhitelistByAllowedValueStmt:      q.getWhitelistByAllowedValueStmt,
		isEmailBlacklistedStmt:              q.isEmailBlacklistedStmt,
		isEmailWhitelistedStmt:              q.isEmailWhitelistedStmt,
		isUserDisabledStmt:                  q.isUserDisabledStmt,
		linkDeviceToUserStmt:                q.linkDeviceToUserStmt,
		updateKYCStatusStmt:                 q.updateKYCStatusStmt,
		updatePublicKeyStmt:                 q.updatePublicKeyStmt,
		updateUserEmailStmt:                 q.updateUserEmailStmt,
		updateUserPasswordStmt:              q.updateUserPasswordStmt,
		updateUserSanitizedEmailStmt:        q.updateUserSanitizedEmailStmt,
		updateUserStatusStmt:                q.updateUserStatusStmt,
		updateUserVerifiedAtStmt:            q.updateUserVerifiedAtStmt,
		updateUsernameStmt:                  q.updateUsernameStmt,
	}
}

// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const blockUsersWithDuplicateEmail = `-- name: BlockUsersWithDuplicateEmail :exec
UPDATE users SET disabled = TRUE, block_reason = 'detected scam: multiple accounts with duplicate email address'
WHERE sanitized_email IN (
        SELECT users.sanitized_email
        FROM users 
        WHERE users.sanitized_email <> '' AND users.sanitized_email IS NOT NULL
        GROUP BY users.sanitized_email
        HAVING count(users.id) > 1 
    )
AND sanitized_email NOT IN (SELECT allowed_value FROM whitelist WHERE allowed_type = 'email')
AND disabled = FALSE
`

func (q *Queries) BlockUsersWithDuplicateEmail(ctx context.Context) error {
	_, err := q.exec(ctx, q.blockUsersWithDuplicateEmailStmt, blockUsersWithDuplicateEmail)
	return err
}

const countAllUsers = `-- name: CountAllUsers :one
SELECT count(id)
FROM users
WHERE verified_at IS NOT NULL
`

func (q *Queries) CountAllUsers(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countAllUsersStmt, countAllUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, username, password, role, sanitized_email)
VALUES ($1, $2, $3, $4, $5) RETURNING id, username, email, password, disabled, verified_at, updated_at, created_at, role, block_reason, sanitized_email, email_hash, kyc_status
`

type CreateUserParams struct {
	Email          string         `json:"email"`
	Username       string         `json:"username"`
	Password       []byte         `json:"password"`
	Role           string         `json:"role"`
	SanitizedEmail sql.NullString `json:"sanitized_email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser,
		arg.Email,
		arg.Username,
		arg.Password,
		arg.Role,
		arg.SanitizedEmail,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Disabled,
		&i.VerifiedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Role,
		&i.BlockReason,
		&i.SanitizedEmail,
		&i.EmailHash,
		&i.KycStatus,
	)
	return i, err
}

const deleteUserByID = `-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteUserByIDStmt, deleteUserByID, id)
	return err
}

const destroyUser = `-- name: DestroyUser :exec
UPDATE users
SET email = 'deleted',
    username = 'deleted',
    password = NULL,
    disabled = TRUE
WHERE id = $1
`

func (q *Queries) DestroyUser(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.destroyUserStmt, destroyUser, userID)
	return err
}

const getKYCStatus = `-- name: GetKYCStatus :one
SELECT kyc_status
FROM users
WHERE id = $1
    LIMIT 1
`

func (q *Queries) GetKYCStatus(ctx context.Context, id uuid.UUID) (sql.NullString, error) {
	row := q.queryRow(ctx, q.getKYCStatusStmt, getKYCStatus, id)
	var kyc_status sql.NullString
	err := row.Scan(&kyc_status)
	return kyc_status, err
}

const getNotSanitizedUsersListDesc = `-- name: GetNotSanitizedUsersListDesc :many
SELECT id, username, email, password, disabled, verified_at, updated_at, created_at, role, block_reason, sanitized_email, email_hash, kyc_status
FROM users
WHERE (sanitized_email IS NULL OR sanitized_email = '')
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetNotSanitizedUsersListDescParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetNotSanitizedUsersListDesc(ctx context.Context, arg GetNotSanitizedUsersListDescParams) ([]User, error) {
	rows, err := q.query(ctx, q.getNotSanitizedUsersListDescStmt, getNotSanitizedUsersListDesc, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Disabled,
			&i.VerifiedAt,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.Role,
			&i.BlockReason,
			&i.SanitizedEmail,
			&i.EmailHash,
			&i.KycStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, email, password, disabled, verified_at, updated_at, created_at, role, block_reason, sanitized_email, email_hash, kyc_status
FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Disabled,
		&i.VerifiedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Role,
		&i.BlockReason,
		&i.SanitizedEmail,
		&i.EmailHash,
		&i.KycStatus,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, email, password, disabled, verified_at, updated_at, created_at, role, block_reason, sanitized_email, email_hash, kyc_status
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.queryRow(ctx, q.getUserByIDStmt, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Disabled,
		&i.VerifiedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Role,
		&i.BlockReason,
		&i.SanitizedEmail,
		&i.EmailHash,
		&i.KycStatus,
	)
	return i, err
}

const getUserBySanitizedEmail = `-- name: GetUserBySanitizedEmail :one
SELECT id, username, email, password, disabled, verified_at, updated_at, created_at, role, block_reason, sanitized_email, email_hash, kyc_status
FROM users
WHERE sanitized_email = $1::text
LIMIT 1
`

func (q *Queries) GetUserBySanitizedEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getUserBySanitizedEmailStmt, getUserBySanitizedEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Disabled,
		&i.VerifiedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Role,
		&i.BlockReason,
		&i.SanitizedEmail,
		&i.EmailHash,
		&i.KycStatus,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, password, disabled, verified_at, updated_at, created_at, role, block_reason, sanitized_email, email_hash, kyc_status
FROM users
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.queryRow(ctx, q.getUserByUsernameStmt, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Disabled,
		&i.VerifiedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Role,
		&i.BlockReason,
		&i.SanitizedEmail,
		&i.EmailHash,
		&i.KycStatus,
	)
	return i, err
}

const getUsersListDesc = `-- name: GetUsersListDesc :many
SELECT id, username, email, password, disabled, verified_at, updated_at, created_at, role, block_reason, sanitized_email, email_hash, kyc_status
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetUsersListDescParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUsersListDesc(ctx context.Context, arg GetUsersListDescParams) ([]User, error) {
	rows, err := q.query(ctx, q.getUsersListDescStmt, getUsersListDesc, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Disabled,
			&i.VerifiedAt,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.Role,
			&i.BlockReason,
			&i.SanitizedEmail,
			&i.EmailHash,
			&i.KycStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVerifiedUsersListDesc = `-- name: GetVerifiedUsersListDesc :many
SELECT id, username, email, password, disabled, verified_at, updated_at, created_at, role, block_reason, sanitized_email, email_hash, kyc_status
FROM users
WHERE verified_at IS NOT NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetVerifiedUsersListDescParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetVerifiedUsersListDesc(ctx context.Context, arg GetVerifiedUsersListDescParams) ([]User, error) {
	rows, err := q.query(ctx, q.getVerifiedUsersListDescStmt, getVerifiedUsersListDesc, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.Disabled,
			&i.VerifiedAt,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.Role,
			&i.BlockReason,
			&i.SanitizedEmail,
			&i.EmailHash,
			&i.KycStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const isUserDisabled = `-- name: IsUserDisabled :one
SELECT disabled
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) IsUserDisabled(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.queryRow(ctx, q.isUserDisabledStmt, isUserDisabled, id)
	var disabled bool
	err := row.Scan(&disabled)
	return disabled, err
}

const updateKYCStatus = `-- name: UpdateKYCStatus :exec
UPDATE users
SET kyc_status = $1::text
WHERE id = $2
`

type UpdateKYCStatusParams struct {
	KycStatus string    `json:"kyc_status"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) UpdateKYCStatus(ctx context.Context, arg UpdateKYCStatusParams) error {
	_, err := q.exec(ctx, q.updateKYCStatusStmt, updateKYCStatus, arg.KycStatus, arg.ID)
	return err
}

const updateUserEmail = `-- name: UpdateUserEmail :exec
UPDATE users
SET email = $2, sanitized_email = $3
WHERE id = $1
`

type UpdateUserEmailParams struct {
	ID             uuid.UUID      `json:"id"`
	Email          string         `json:"email"`
	SanitizedEmail sql.NullString `json:"sanitized_email"`
}

func (q *Queries) UpdateUserEmail(ctx context.Context, arg UpdateUserEmailParams) error {
	_, err := q.exec(ctx, q.updateUserEmailStmt, updateUserEmail, arg.ID, arg.Email, arg.SanitizedEmail)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1
`

type UpdateUserPasswordParams struct {
	ID       uuid.UUID `json:"id"`
	Password []byte    `json:"password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.exec(ctx, q.updateUserPasswordStmt, updateUserPassword, arg.ID, arg.Password)
	return err
}

const updateUserSanitizedEmail = `-- name: UpdateUserSanitizedEmail :exec
UPDATE users
SET sanitized_email = $1::text
WHERE id = $2
`

type UpdateUserSanitizedEmailParams struct {
	SanitizedEmail string    `json:"sanitized_email"`
	ID             uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserSanitizedEmail(ctx context.Context, arg UpdateUserSanitizedEmailParams) error {
	_, err := q.exec(ctx, q.updateUserSanitizedEmailStmt, updateUserSanitizedEmail, arg.SanitizedEmail, arg.ID)
	return err
}

const updateUserStatus = `-- name: UpdateUserStatus :exec
UPDATE users
SET disabled = $2, block_reason = $3
WHERE id = $1
`

type UpdateUserStatusParams struct {
	ID          uuid.UUID      `json:"id"`
	Disabled    bool           `json:"disabled"`
	BlockReason sql.NullString `json:"block_reason"`
}

func (q *Queries) UpdateUserStatus(ctx context.Context, arg UpdateUserStatusParams) error {
	_, err := q.exec(ctx, q.updateUserStatusStmt, updateUserStatus, arg.ID, arg.Disabled, arg.BlockReason)
	return err
}

const updateUserVerifiedAt = `-- name: UpdateUserVerifiedAt :exec
UPDATE users
SET verified_at = $1
WHERE id = $2
`

type UpdateUserVerifiedAtParams struct {
	VerifiedAt sql.NullTime `json:"verified_at"`
	UserID     uuid.UUID    `json:"user_id"`
}

func (q *Queries) UpdateUserVerifiedAt(ctx context.Context, arg UpdateUserVerifiedAtParams) error {
	_, err := q.exec(ctx, q.updateUserVerifiedAtStmt, updateUserVerifiedAt, arg.VerifiedAt, arg.UserID)
	return err
}

const updateUsername = `-- name: UpdateUsername :exec
UPDATE users
SET username = $2
WHERE id = $1
`

type UpdateUsernameParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (q *Queries) UpdateUsername(ctx context.Context, arg UpdateUsernameParams) error {
	_, err := q.exec(ctx, q.updateUsernameStmt, updateUsername, arg.ID, arg.Username)
	return err
}

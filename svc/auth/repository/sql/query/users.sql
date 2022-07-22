-- name: GetUsersListDesc :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetVerifiedUsersListDesc :many
SELECT *
FROM users
WHERE verified_at IS NOT NULL
AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllUsers :one
SELECT count(id)
FROM users
WHERE verified_at IS NOT NULL AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1 
AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = @email
AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserBySanitizedEmail :one
SELECT *
FROM users
WHERE sanitized_email = @email::text
AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
AND deleted_at IS NULL
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, username, password, role, sanitized_email, email_hash)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateUserEmail :exec
UPDATE users
SET email = $2, sanitized_email = $3, email_hash = $4
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUsername :exec
UPDATE users
SET username = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserStatus :exec
UPDATE users
SET disabled = $2, block_reason = $3
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserVerifiedAt :exec
UPDATE users
SET verified_at = @verified_at
WHERE id = @user_id AND deleted_at IS NULL;

-- name: DestroyUser :exec
UPDATE users
SET email = email_hash,
    username = email_hash,
    password = NULL,
    disabled = TRUE,
    sanitized_email = NULL,
    public_key = NULL,
    deleted_at = NOW()
WHERE id = @user_id AND deleted_at IS NULL;

-- name: IsUserDisabled :one
SELECT disabled
FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;


-- name: BlockUsersWithDuplicateEmail :exec 
UPDATE users SET disabled = TRUE, block_reason = 'detected scam: multiple accounts with duplicate email address'
WHERE sanitized_email IN (
        SELECT users.sanitized_email
        FROM users 
        WHERE users.sanitized_email <> '' 
            AND users.sanitized_email IS NOT NULL 
            AND users.verified_at IS NOT NULL
        GROUP BY users.sanitized_email
        HAVING count(users.id) > 1 
    )
AND sanitized_email NOT IN (SELECT allowed_value FROM whitelist WHERE allowed_type = 'email')
AND disabled = FALSE
AND deleted_at IS NULL;

-- name: GetNotSanitizedUsersListDesc :many
SELECT *
FROM users
WHERE (sanitized_email IS NULL OR sanitized_email = '')
AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserSanitizedEmail :exec
UPDATE users
SET sanitized_email = @sanitized_email::text
WHERE id = @id AND deleted_at IS NULL;

-- name: UpdateKYCStatus :exec
UPDATE users
SET kyc_status = @kyc_status::text
WHERE id = @id AND deleted_at IS NULL;

-- name: GetKYCStatus :one
SELECT kyc_status::text
FROM users
WHERE id = $1 AND deleted_at IS NULL
    LIMIT 1;

-- name: GetUsernameByID :one
SELECT username 
FROM users
WHERE id = @id AND deleted_at IS NULL
    LIMIT 1;  

-- name: UpdatePublicKey :exec
UPDATE users
SET public_key = @public_key::text
WHERE id = @id AND deleted_at IS NULL;

-- name: GetPublicKey :one
SELECT public_key
FROM users
WHERE id = @id AND deleted_at IS NULL
    LIMIT 1;

-- name: UpdateUserRole :exec
UPDATE users
SET role = @role
WHERE id = @id AND deleted_at IS NULL;
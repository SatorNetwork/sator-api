-- name: GetUsersListDesc :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetVerifiedUsersListDesc :many
SELECT *
FROM users
WHERE verified_at IS NOT NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllUsers :one
SELECT count(id)
FROM users
WHERE verified_at IS NOT NULL;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;
-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;
-- name: CreateUser :one
INSERT INTO users (email, username, password, role)
VALUES ($1, $2, $3, $4) RETURNING *;
-- name: UpdateUserEmail :exec
UPDATE users
SET email = $2
WHERE id = $1;
-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1;
-- name: UpdateUserStatus :exec
UPDATE users
SET disabled = $2
WHERE id = $1;
-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;
-- name: UpdateUserVerifiedAt :exec
UPDATE users
SET verified_at = @verified_at
WHERE id = @user_id;
-- name: DestroyUser :exec
UPDATE users
SET email = 'deleted',
    username = 'deleted',
    password = NULL,
    disabled = TRUE
WHERE id = @user_id;
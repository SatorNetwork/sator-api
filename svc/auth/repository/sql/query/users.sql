-- name: GetUsersListDesc :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
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
INSERT INTO users (email, username, password)
VALUES ($1, $2, $3) RETURNING *;
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
SET verified_at = $1;
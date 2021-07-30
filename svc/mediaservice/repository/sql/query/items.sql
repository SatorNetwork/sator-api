-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1
    LIMIT 1;
-- name: GetItemsList :many
SELECT * FROM items
LIMIT $1 OFFSET $2;
-- name: CreateItem :one
INSERT INTO items (id, filename, filepath, relation_type, relation_id)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;
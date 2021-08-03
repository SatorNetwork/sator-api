-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1
    LIMIT 1;
-- name: GetItemsList :many
SELECT * FROM items
LIMIT $1 OFFSET $2;
-- name: AddItem :one
INSERT INTO items (id, filename, filepath)
VALUES ($1, $2, $3)
    RETURNING *;
-- name: DeleteItemByID :exec
DELETE FROM items
WHERE id = @id;
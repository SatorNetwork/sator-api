-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1
    LIMIT 1;
-- name: GetItemsList :many
SELECT * FROM items
LIMIT $1 OFFSET $2;
-- name: AddItem :one
INSERT INTO items (id, file_name, file_path, file_url)
VALUES ($1, $2, $3, $4)
    RETURNING *;
-- name: DeleteItemByID :exec
DELETE FROM items
WHERE id = @id;
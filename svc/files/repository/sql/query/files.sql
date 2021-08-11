-- name: GetFileByID :one
SELECT * FROM files
WHERE id = $1
    LIMIT 1;
-- name: GetFilesList :many
SELECT * FROM files
LIMIT $1 OFFSET $2;
-- name: AddFile :one
INSERT INTO files (id, file_name, file_path, file_url)
VALUES ($1, $2, $3, $4)
    RETURNING *;
-- name: DeleteFileByID :exec
DELETE FROM files
WHERE id = @id;
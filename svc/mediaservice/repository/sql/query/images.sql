-- name: GetImageByID :one
SELECT * FROM images
WHERE id = $1
    LIMIT 1;
-- name: GetImagesList :many
SELECT * FROM images
LIMIT $1 OFFSET $2;
-- name: AddImage :one
INSERT INTO images (id, file_name, file_path, file_url)
VALUES ($1, $2, $3, $4)
    RETURNING *;
-- name: DeleteImageByID :exec
DELETE FROM images
WHERE id = @id;
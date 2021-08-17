-- name: GetShows :many
SELECT *
FROM shows
ORDER BY has_new_episode DESC,
    updated_at DESC,
    created_at DESC
LIMIT $1 OFFSET $2;
-- name: GetShowByID :one
SELECT *
FROM shows
WHERE id = $1;
-- name: GetShowsByCategory :many
SELECT *
FROM shows
WHERE category = $1
ORDER BY has_new_episode DESC,
         updated_at DESC,
         created_at DESC
    LIMIT $2 OFFSET $3;
-- name: AddShow :one
INSERT INTO shows (
    title,
    cover,
    has_new_episode,
    category,
    description
  )
VALUES (
           @title,
           @cover,
           @has_new_episode,
           @category,
           @description
) RETURNING *;
-- name: UpdateShow :exec
UPDATE shows
SET title = @title,
    cover = @cover,
    has_new_episode = @has_new_episode,
    category = @category,
    description = @description
WHERE id = @id;
-- name: DeleteShowByID :exec
DELETE FROM shows
WHERE id = @id;
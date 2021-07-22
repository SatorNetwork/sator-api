-- name: GetShowsPaginated :many
SELECT *
FROM shows
LEFT JOIN shows_to_category
ON id = show_id
ORDER BY has_new_episode DESC,
    updated_at DESC,
    created_at DESC
LIMIT $1 OFFSET $2;
-- name: GetShowByID :one
SELECT *
FROM shows
LEFT JOIN shows_to_category
ON id = show_id
WHERE id = $1;
-- name: AddShow :one
WITH inserted_shows AS(

INSERT INTO shows (
    title,
    cover,
    has_new_episode,
    description
  )
VALUES (
           @title,
           @cover,
           @has_new_episode,
           @description
       ) RETURNING * )
SELECT * FROM inserted_shows
LEFT JOIN shows_to_category
ON inserted_shows.id = shows_to_category.show_id;
-- name: UpdateShow :exec
UPDATE shows
SET title = @title,
    cover = @cover,
    has_new_episode = @has_new_episode,
    description = @description
WHERE id = @id;
-- name: DeleteShowByID :exec
DELETE FROM shows
WHERE id = @id;
-- name: GetShows :many
SELECT *
FROM shows
         LEFT JOIN shows_to_category
                   ON id = show_id
ORDER BY has_new_episode DESC,
         updated_at DESC,
         created_at DESC;
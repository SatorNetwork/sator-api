-- name: GetShows :many
WITH show_claps_sum AS (
    SELECT 
        COUNT(*) AS claps,
        show_id
    FROM show_claps
    GROUP BY show_id  
)
SELECT shows.*, show_claps_sum.claps AS claps
FROM shows
LEFT JOIN show_claps_sum ON show_claps_sum.show_id = shows.id
ORDER BY has_new_episode DESC,
    updated_at DESC,
    created_at DESC
LIMIT $1 OFFSET $2;
-- name: GetShowByID :one
WITH show_claps_sum AS (
    SELECT 
        COUNT(*) AS claps,
        show_id
    FROM show_claps
    WHERE id = @id
    GROUP BY show_id  
)
SELECT 
    shows.*,
    show_claps_sum.claps as claps
FROM shows
LEFT JOIN show_claps_sum ON show_claps_sum.show_id = shows.id
WHERE shows.id = @id;
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
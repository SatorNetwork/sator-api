-- name: GetShows :many
SELECT *
FROM shows
WHERE archived = FALSE
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
    WHERE show_id = @id
    GROUP BY show_id  
)
SELECT 
    shows.*,
    COALESCE(show_claps_sum.claps, 0) as claps
FROM shows
LEFT JOIN show_claps_sum ON show_claps_sum.show_id = shows.id
WHERE shows.id = @id AND shows.archived = FALSE;
-- name: GetShowsByCategory :many
WITH filtred_shows AS (
    SELECT show_id
    FROM shows_to_category
    WHERE category_id = $1 AND disabled = false
    GROUP BY show_id
    )
SELECT
       shows.*,
       coalesce (filtred_shows.show_id, 0) as show_id
FROM shows
INNER JOIN filtred_shows ON shows.id = filtred_shows.show_id
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
    description,
    realms_title,
    realms_subtitle,
    watch
  )
VALUES (
           @title,
           @cover,
           @has_new_episode,
           @category,
           @description,
           @realms_title,
           @realms_subtitle,
           @watch
) RETURNING *;
-- name: UpdateShow :exec
UPDATE shows
SET title = @title,
    cover = @cover,
    has_new_episode = @has_new_episode,
    category = @category,
    description = @description,
    realms_title = @realms_title,
    realms_subtitle = @realms_subtitle,
    watch = @watch
WHERE id = @id;
-- name: DeleteShowByID :exec
DELETE FROM shows
WHERE id = @id;
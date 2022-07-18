-- name: GetPublishedShows :many
SELECT *
FROM shows
WHERE status = 'published'::shows_status_type
AND shows.deleted_at IS NULL
ORDER BY has_new_episode DESC,
    updated_at DESC,
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAllShows :many
SELECT *
FROM shows
WHERE shows.deleted_at IS NULL
ORDER BY has_new_episode DESC,
    updated_at DESC,
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetShowsByOldCategory :many
SELECT *
FROM shows
WHERE status = 'published'::shows_status_type
AND category = @category::varchar
AND shows.deleted_at IS NULL
ORDER BY has_new_episode DESC,
    updated_at DESC,
    created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetShowByID :one
SELECT *
FROM shows
WHERE shows.id = @id AND shows.deleted_at IS NULL;

-- name: GetPublishedShowByID :one
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
WHERE shows.id = @id AND shows.status = 'published'::shows_status_type AND shows.deleted_at IS NULL;

-- name: GetShowsByCategory :many
SELECT * FROM shows
WHERE id IN(
        SELECT DISTINCT show_id FROM shows_to_categories
              JOIN show_categories ON show_categories.id = shows_to_categories.category_id
        WHERE show_categories.disabled = FALSE
          AND show_categories.id = @category_id)
AND status = 'published'::shows_status_type
AND shows.deleted_at IS NULL
ORDER BY has_new_episode DESC,
         updated_at DESC,
         created_at DESC
    LIMIT @limit_val OFFSET @offset_val;

-- name: AddShow :one
INSERT INTO shows (
    title,
    cover,
    has_new_episode,
    category,
    description,
    realms_title,
    realms_subtitle,
    watch,
    status
  )
VALUES (
    @title,
    @cover,
    @has_new_episode,
    @category,
    @description,
    @realms_title,
    @realms_subtitle,
    @watch,
    @status::shows_status_type
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
    watch = @watch,
    status = @status::shows_status_type
WHERE id = @id AND shows.deleted_at IS NULL;

-- name: DeleteShowByID :exec
UPDATE shows
SET deleted_at = NOW()
WHERE id = @id AND shows.deleted_at IS NULL;

-- name: GetShowsByTitle :many
SELECT * FROM shows
WHERE title = @title AND shows.deleted_at IS NULL;

-- name: GetShowsByStatus :many
SELECT *
FROM shows
WHERE status = @status::shows_status_type
AND shows.deleted_at IS NULL
LIMIT @limit_val OFFSET @offset_val;
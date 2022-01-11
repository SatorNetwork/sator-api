-- name: GetShowCategories :many
SELECT *
FROM shows_categories
WHERE disabled = FALSE
ORDER BY sort DESC
    LIMIT $1 OFFSET $2;
-- name: GetShowCategoryByID :one
SELECT *
FROM shows_categories
WHERE id = $1;

-- name: AddShowCategory :one
INSERT INTO shows_categories (
    sort,
    title,
    disabled
)
VALUES (
           @sort,
           @title,
           @disabled
       ) RETURNING *;

-- name: UpdateShowCategory :exec
UPDATE shows_categories
SET sort = @sort,
    title = @title,
    disabled = @disabled
WHERE id = @id;

-- name: DeleteShowCategoryByID :exec
DELETE FROM shows_categories
WHERE id = @id;
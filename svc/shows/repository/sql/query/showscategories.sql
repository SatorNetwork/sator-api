-- name: GetShowCategoryByID :one
SELECT *
FROM shows_categories
WHERE id = $1;
-- name: AddShowCategory :one
INSERT INTO shows_categories (
    category_name,
    title,
    disabled
)
VALUES (
           @category_name,
           @title,
           @disabled
       ) RETURNING *;
-- name: UpdateShowCategory :exec
UPDATE shows_categories
SET category_name = @category_name,
    title = @title,
    disabled = @disabled
WHERE id = @id;
-- name: DeleteShowCategoryByID :exec
DELETE FROM shows_categories
WHERE id = @id;
-- name: GetShowCategories :many
SELECT *
FROM shows_categories;
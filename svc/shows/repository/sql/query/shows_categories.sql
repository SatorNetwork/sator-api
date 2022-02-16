-- name: GetShowCategories :many
SELECT *
FROM show_categories
WHERE disabled = FALSE
ORDER BY sort ASC
LIMIT $1 OFFSET $2;

-- name: GetShowCategoriesWithDisabled :many
SELECT *
FROM show_categories
ORDER BY sort ASC
LIMIT $1 OFFSET $2;

-- name: GetShowCategoryByID :one
SELECT *
FROM show_categories
WHERE id = $1;

-- name: AddShowCategory :one
INSERT INTO show_categories (
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
UPDATE show_categories
SET sort = @sort,
    title = @title,
    disabled = @disabled
WHERE id = @id;

-- name: DeleteShowCategoryByID :exec
DELETE FROM show_categories
WHERE id = @id;

-- name: AddShowToCategory :one
INSERT INTO shows_to_categories (
    category_id,
    show_id
    )
VALUES (
           @category_id,
           @show_id
       ) RETURNING *;

-- name: DeleteShowToCategoryByShowID :exec
DELETE FROM shows_to_categories
WHERE show_id = @show_id;

-- name: GetCategoriesByShowID :many
SELECT category_id
FROM shows_to_categories
WHERE show_id = $1;

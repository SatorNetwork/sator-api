-- name: AddShowToCategory :exec
INSERT INTO shows_to_category (
    category_id,
    show_id
)
VALUES (
           @category_id,
           @show_id
       );
-- name: UpdateShowToCategory :exec
UPDATE shows_to_category
SET category_id = @category_id,
    show_id = @show_id
WHERE category_id = @category_id AND show_id = @show_id;
-- name: DeleteShowToCategory :exec
DELETE FROM shows_to_category
WHERE category_id = @category_id AND show_id = @show_id;
-- name: DeleteShowToCategoryByShowID :exec
DELETE FROM shows_to_category
WHERE show_id = @show_id;
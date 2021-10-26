-- name: GetNFTCategoryByID :one
SELECT * FROM nft_categories
WHERE id = @id
LIMIT 1;

-- name: GetMainNFTCategory :one
SELECT * FROM nft_categories
WHERE main = TRUE
LIMIT 1;

-- name: GetNFTCategoriesList :many
SELECT * FROM nft_categories
ORDER BY sort ASC;

-- name: AddNFTCategory :one
INSERT INTO nft_categories (title, sort)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateNFTCategory :one
UPDATE nft_categories SET title = @title, main = @main
WHERE id = @id
RETURNING *;

-- name: ResetMainNFTCategory :exec
UPDATE nft_categories SET main = FALSE
WHERE main = TRUE;

-- name: DeleteNFTCategoryByID :exec
DELETE FROM nft_categories
WHERE id = @id;
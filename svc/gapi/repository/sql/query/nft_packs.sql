-- name: GetNFTPacksList :many
SELECT * FROM unity_game_nft_packs WHERE deleted_at IS NULL;

-- name: AddNFTPack :one
INSERT INTO unity_game_nft_packs ( name, drop_chances, price) VALUES ($1, $2, $3) RETURNING *;

-- name: GetNFTPack :one
SELECT * FROM unity_game_nft_packs WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateNFTPack :one
UPDATE unity_game_nft_packs SET drop_chances = $1, price = $2, name = $3 WHERE id = $4 RETURNING *;

-- name: DeleteNFTPack :exec
DELETE FROM unity_game_nft_packs WHERE id = $1;

-- name: SoftDeleteNFTPack :exec
UPDATE unity_game_nft_packs SET deleted_at = now() WHERE id = $1;
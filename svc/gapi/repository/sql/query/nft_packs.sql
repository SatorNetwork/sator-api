-- name: GetNFTPacksList :many
SELECT * FROM unity_game_nft_packs WHERE deleted_at IS NULL;

-- name: AddNFTPack :exec
INSERT INTO unity_game_nft_packs (id, drop_chances, price) VALUES ($1, $2, $3);

-- name: GetNFTPack :one
SELECT * FROM unity_game_nft_packs WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateNFTPack :exec
UPDATE unity_game_nft_packs SET drop_chances = $1, price = $2 WHERE id = $3;

-- name: DeleteNFTPack :exec
DELETE FROM unity_game_nft_packs WHERE id = $1;

-- name: SoftDeleteNFTPack :exec
UPDATE unity_game_nft_packs SET deleted_at = now() WHERE id = $1;
-- name: AddNFT :one
INSERT INTO unity_game_nfts (id, nft_type, allowed_levels)  
VALUES ($1, $2, $3) RETURNING *;

-- name: GetNFT :one
SELECT * FROM unity_game_nfts WHERE id = $1;

-- name: GetNFTs :many
SELECT * FROM unity_game_nfts LIMIT $1 OFFSET $2;

-- name: GetNFTsByTypeAndLevel :many
SELECT * FROM unity_game_nfts WHERE nft_type = $1 AND allowed_levels @> $2;

-- name: UpdateNFT :exec
UPDATE unity_game_nfts SET nft_type = $1, allowed_levels = $2 WHERE id = $3;

-- name: DeleteNFT :exec
DELETE FROM unity_game_nfts WHERE id = $1;

-- name: SoftDeleteNFT :exec
UPDATE unity_game_nfts SET deleted_at = now() WHERE id = $1;
-- name: AddNFT :one
INSERT INTO unity_game_nfts (id, user_id, nft_type, max_level)  
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetNFT :one
SELECT * FROM unity_game_nfts WHERE id = $1;

-- name: GetUserNFT :one
SELECT * FROM unity_game_nfts WHERE id = $1 AND user_id = $2;

-- name: GetUserNFTs :many
SELECT * FROM unity_game_nfts 
WHERE user_id = $1 
AND deleted_at IS NULL 
AND crafted_nft_id IS NULL;

-- name: GetUserNFTByIDs :many
SELECT * FROM unity_game_nfts 
WHERE user_id = @user_id
AND id = ANY(@ids::VARCHAR[])
AND deleted_at IS NULL 
AND crafted_nft_id IS NULL;

-- name: DeleteNFT :exec
UPDATE unity_game_nfts SET deleted_at = now() WHERE id = $1;

-- name: CraftNFTs :exec
UPDATE unity_game_nfts SET crafted_nft_id = @crafted_nft_id, deleted_at = now()
WHERE user_id = @user_id AND id = ANY(@nft_ids::VARCHAR[]);
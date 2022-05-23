-- name: LinkNFTToPlayer :exec
INSERT INTO unity_game_player_nfts (user_id, nft_id) VALUES ($1, $2);

-- name: GetNFTsByPlayer :many
SELECT unity_game_nfts.* 
FROM unity_game_player_nfts 
JOIN unity_game_nfts ON unity_game_player_nfts.nft_id = unity_game_nfts.id AND unity_game_nfts.deleted_at IS NULL
WHERE user_id = $1
AND crafted_nft_id IS NULL;

-- name: UnlinkNFTFromPlayer :exec
DELETE FROM unity_game_player_nfts WHERE user_id = $1 AND nft_id = $2;

-- name: CraftNFTs :exec
UPDATE unity_game_player_nfts SET crafted_nft_id = @crafted_nft_id
WHERE user_id = @user_id AND nft_id = ANY(@nft_ids::VARCHAR[]);
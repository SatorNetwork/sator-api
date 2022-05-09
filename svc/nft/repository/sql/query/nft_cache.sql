-- name: AddNFTToCache :exec
INSERT INTO nft_cache (mint_addr, metadata) 
VALUES ($1, $2);

-- name: GetNFTFromCache :one
SELECT * FROM nft_cache
WHERE mint_addr = $1;
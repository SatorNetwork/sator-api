-- name: GetNFTItemByID :one
SELECT * FROM nft_items
WHERE id = $1
LIMIT 1;

-- name: GetNFTItemsList :many
SELECT * FROM nft_items
ORDER BY updated_at DESC, created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetNFTItemsListByOwnerID :many
SELECT * FROM nft_items
WHERE owner_id = @owner_id
ORDER BY updated_at DESC, created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetNFTItemsListByRelationID :many
SELECT nft_items.* FROM nft_items
JOIN nft_relations ON nft_relations.nft_item_id =  nft_items.id
WHERE nft_relations.relation_id = @relation_id
AND nft_items.owner_id IS NULL
ORDER BY nft_items.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: AddNFTItem :one
INSERT INTO nft_items (name, description, cover, supply, buy_now_price, token_uri)
VALUES (@name, @description, @cover, @supply, @buy_now_price, @token_uri)
RETURNING *;

-- name: UpdateNFTItemOwner :exec
UPDATE nft_items SET owner_id = @owner_id
WHERE id = @id;
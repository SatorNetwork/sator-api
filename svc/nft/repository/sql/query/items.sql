-- name: GetNFTItemByID :one
WITH minted_nft_items AS (
    SELECT COUNT(user_id)::INT as minted, nft_item_id
    FROM nft_owners
    WHERE nft_owners.nft_item_id = @id
    GROUP BY nft_owners.nft_item_id
)
SELECT nft_items.*,
    coalesce(minted_nft_items.minted, 0) AS minted
FROM nft_items
LEFT JOIN minted_nft_items ON minted_nft_items.nft_item_id = nft_items.id
WHERE id = @id
AND nft_items.supply > 0
LIMIT 1;

-- name: DoesUserOwnNFT :one
SELECT count(*) > 0
FROM nft_owners
WHERE user_id = @user_id
AND nft_item_id = @nft_item_id;

-- name: GetNFTItemsList :many
WITH minted_nfts AS (
    SELECT nft_item_id, COUNT(user_id)::INT AS minted
    FROM nft_owners
    GROUP BY nft_item_id
)
SELECT nft_items.*, minted_nfts.minted as minted
FROM nft_items
    LEFT JOIN minted_nfts ON minted_nfts.nft_item_id = nft_items.id
WHERE nft_items.supply > 0
ORDER BY nft_items.updated_at DESC, nft_items.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetNFTItemsListByOwnerID :many
SELECT *
FROM nft_items
WHERE nft_items.id = ANY(SELECT DISTINCT nft_owners.nft_item_id
                        FROM nft_owners
                        WHERE nft_owners.user_id = @owner_id)
ORDER BY nft_items.updated_at DESC, nft_items.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetNFTItemsListByRelationID :many
WITH minted_nfts AS (
    SELECT nft_item_id, COUNT(user_id)::INT AS minted
    FROM nft_owners
    GROUP BY nft_item_id
)
SELECT nft_items.*, minted_nfts.minted as minted
FROM nft_items
    LEFT JOIN minted_nfts ON minted_nfts.nft_item_id = nft_items.id
WHERE nft_items.id = ANY(SELECT DISTINCT nft_relations.nft_item_id 
                        FROM nft_relations 
                        WHERE nft_relations.relation_id = @relation_id)
AND (nft_items.supply > minted_nfts.minted OR minted_nfts.minted IS NULL)
AND nft_items.supply > 0
ORDER BY nft_items.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetAllNFTItems :many
WITH minted_nfts AS (
    SELECT nft_item_id, COUNT(user_id)::INT AS minted
    FROM nft_owners
    GROUP BY nft_item_id
)
SELECT nft_items.*, minted_nfts.minted as minted
FROM nft_items
    LEFT JOIN minted_nfts ON minted_nfts.nft_item_id = nft_items.id
WHERE nft_items.supply > COALESCE(minted_nfts.minted, 0) 
AND nft_items.supply > 0
ORDER BY nft_items.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: AddNFTItem :one
INSERT INTO nft_items (name, description, cover, supply, buy_now_price, token_uri)
VALUES (@name, @description, @cover, @supply, @buy_now_price, @token_uri)
RETURNING *;

-- name: AddNFTItemOwner :exec
INSERT INTO nft_owners (nft_item_id, user_id)
VALUES (@nft_item_id, @user_id);

-- name: DeleteNFTItemByID :exec
DELETE FROM nft_items
WHERE id = @id;

-- name: UpdateNFTItem :exec
UPDATE nft_items
SET name = @name,
    description = @description,
    cover = @cover,
    supply = @supply,
    buy_now_price = @buy_now_price,
    token_uri = @token_uri
WHERE id = @id;
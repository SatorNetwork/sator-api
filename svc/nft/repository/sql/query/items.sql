-- name: GetNFTItemByID :one
SELECT nft_items.*,
    (SELECT COUNT(user_id)::INT
    FROM nft_owners
    WHERE nft_owners.nft_item_id = @id
    GROUP BY nft_owners.nft_item_id) AS minted
FROM nft_items
WHERE id = @id
LIMIT 1;

-- name: GetNFTItemsList :many
WITH minted_nfts AS (
    SELECT nft_item_id, COUNT(user_id)::INT AS minted
    FROM nft_owners
    GROUP BY nft_item_id
)
SELECT nft_items.*
FROM nft_items
    LEFT JOIN minted_nfts ON minted_nfts.nft_item_id = nft_items.id
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
SELECT nft_items.*
FROM nft_items
    LEFT JOIN minted_nfts ON minted_nfts.nft_item_id = nft_items.id
    JOIN nft_relations ON nft_relations.nft_item_id = nft_items.id
WHERE nft_items.id = ANY(SELECT DISTINCT nft_relations.nft_item_id 
                        FROM nft_relations 
                        WHERE nft_relations.relation_id = @relation_id)
AND nft_items.supply > minted_nfts.minted
ORDER BY nft_items.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: AddNFTItem :one
INSERT INTO nft_items (name, description, cover, supply, buy_now_price, token_uri)
VALUES (@name, @description, @cover, @supply, @buy_now_price, @token_uri)
RETURNING *;

-- name: AddNFTItemOwner :exec
INSERT INTO nft_owners (nft_item_id, user_id)
VALUES (@nft_item_id, @user_id);
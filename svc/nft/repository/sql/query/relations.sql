-- name: AddNFTRelation :exec
INSERT INTO nft_relations (nft_item_id, relation_id)
VALUES (@nft_item_id, @relation_id);

-- name: DeleteNFTRelation :exec
DELETE FROM nft_relations
WHERE nft_item_id = @nft_item_id
AND relation_id = @relation_id;

-- name: DoesRelationIDHasRelationNFT :one
WITH minted_nft_items AS (
    SELECT COUNT(user_id)::INT as minted, nft_item_id
    FROM nft_owners
    GROUP BY nft_owners.nft_item_id
)
SELECT EXISTS(
    SELECT nft_relations.nft_item_id
    FROM nft_relations
        LEFT JOIN minted_nft_items ON minted_nft_items.nft_item_id = nft_relations.nft_item_id
    WHERE relation_id = $1
);

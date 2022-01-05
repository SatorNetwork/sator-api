-- name: AddNFTRelation :exec
INSERT INTO nft_relations (nft_item_id, relation_id)
VALUES (@nft_item_id, @relation_id);

-- name: DeleteNFTRelation :exec
DELETE FROM nft_relations
WHERE nft_item_id = @nft_item_id
AND relation_id = @relation_id;

-- name: DoesRelationIDHasRelationNFT :one
WITH minted_nfts AS (
    SELECT nft_item_id, COUNT(user_id)::INT AS minted
    FROM nft_owners
    GROUP BY nft_item_id
)
SELECT EXISTS(
               SELECT nft_relations.nft_item_id
               FROM nft_relations
               WHERE relation_id = $1 AND nft_item_id IN (
                   SELECT nft_items.id
                   FROM nft_items
                            LEFT JOIN minted_nfts ON minted_nfts.nft_item_id = nft_items.id
                   WHERE nft_items.supply > COALESCE (minted_nfts.minted, 0)
               )
           );

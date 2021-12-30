-- name: AddNFTRelation :exec
INSERT INTO nft_relations (nft_item_id, relation_id)
VALUES (@nft_item_id, @relation_id);

-- name: DeleteNFTRelation :exec
DELETE FROM nft_relations
WHERE nft_item_id = @nft_item_id
AND relation_id = @relation_id;

-- name: DoesRelationIDHasRelationNFT :one
SELECT EXISTS(
    SELECT nft_item_id FROM nft_relations
    WHERE relation_id = $1
);
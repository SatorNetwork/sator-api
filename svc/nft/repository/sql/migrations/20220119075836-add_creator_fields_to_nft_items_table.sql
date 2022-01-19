-- +migrate Up
ALTER TABLE nft_items
    ADD COLUMN creator_address VARCHAR DEFAULT NULL,
    ADD COLUMN creator_share int DEFAULT 0;

-- +migrate Down 
ALTER TABLE nft_items 
	DROP COLUMN creator_address, 
	DROP COLUMN creator_share;
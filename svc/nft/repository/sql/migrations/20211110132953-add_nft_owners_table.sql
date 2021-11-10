-- +migrate Up
CREATE TABLE IF NOT EXISTS nft_owners (
    nft_item_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
-- +migrate Down
DROP TABLE IF EXISTS nft_owners;
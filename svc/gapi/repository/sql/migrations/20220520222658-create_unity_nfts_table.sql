
-- +migrate Up
CREATE TABLE IF NOT EXISTS unity_game_nfts (
    id VARCHAR PRIMARY KEY,
    user_id uuid NOT NULL,
    nft_type VARCHAR NOT NULL,
    max_level INT NOT NULL DEFAULT 1,
    crafted_nft_id VARCHAR DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS unity_game_nfts;

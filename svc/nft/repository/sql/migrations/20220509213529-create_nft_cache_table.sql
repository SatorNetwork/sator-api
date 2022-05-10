
-- +migrate Up
CREATE TABLE IF NOT EXISTS nft_cache (
    mint_addr VARCHAR PRIMARY KEY,
    metadata bytea NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS nft_cache;

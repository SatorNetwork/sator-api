
-- +migrate Up
CREATE TABLE IF NOT EXISTS unity_game_player_nfts (
    user_id uuid,
    nft_id VARCHAR NOT NULL,
    crafted_nft_id VARCHAR DEFAULT NULL,
    PRIMARY KEY (user_id, nft_id)
);

-- +migrate Down
DROP TABLE IF EXISTS unity_game_player_nfts;

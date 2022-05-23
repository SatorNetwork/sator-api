-- +migrate Up
ALTER TABLE unity_game_nft_packs
    ADD COLUMN name VARCHAR NOT NULL;

-- +migrate Down
ALTER TABLE unity_game_nft_packs DROP COLUMN name;
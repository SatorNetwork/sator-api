
-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION unity_game_nft_packs_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS unity_game_nft_packs (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    drop_chances bytea NOT NULL,
    price DOUBLE PRECISION NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL
);
CREATE INDEX IF NOT EXISTS index_unity_game_nft_packs_on_deleted_at ON unity_game_nft_packs (deleted_at);

-- +migrate Down
CREATE TRIGGER update_unity_game_nft_packs_modtime BEFORE
UPDATE ON unity_game_nft_packs FOR EACH ROW EXECUTE PROCEDURE unity_game_nft_packs_update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_unity_game_nft_packs_modtime ON unity_game_nft_packs;
DROP TABLE IF EXISTS unity_game_nft_packs;
DROP FUNCTION IF EXISTS unity_game_nft_packs_update_updated_at_column();

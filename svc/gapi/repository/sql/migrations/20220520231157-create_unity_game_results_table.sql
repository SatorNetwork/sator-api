
-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION unity_game_results_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS unity_game_results (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    nft_id VARCHAR NOT NULL,
    complexity VARCHAR NOT NULL,
    is_training BOOLEAN NOT NULL DEFAULT false,
    blocks_done INTEGER NOT NULL DEFAULT 0,
    finished_at TIMESTAMP DEFAULT NULL,
    rewards DOUBLE PRECISION NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TRIGGER update_unity_game_results_modtime BEFORE
UPDATE ON unity_game_results FOR EACH ROW EXECUTE PROCEDURE unity_game_results_update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_unity_game_results_modtime ON unity_game_results;
DROP TABLE IF EXISTS unity_game_results;
DROP FUNCTION IF EXISTS unity_game_results_update_updated_at_column();

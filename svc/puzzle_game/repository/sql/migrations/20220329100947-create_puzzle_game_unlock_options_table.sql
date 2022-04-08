-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION puzzle_game_unlock_options_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS puzzle_game_unlock_options (
    id VARCHAR PRIMARY KEY,
    steps INTEGER NOT NULL DEFAULT 0,
    amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TRIGGER update_puzzle_game_unlock_options_modtime BEFORE
    UPDATE ON puzzle_game_unlock_options FOR EACH ROW EXECUTE PROCEDURE puzzle_game_unlock_options_update_updated_at_column();
    
-- +migrate Down
DROP TRIGGER IF EXISTS update_puzzle_game_unlock_options_modtime ON puzzle_game_unlock_options;
DROP TABLE IF EXISTS puzzle_game_unlock_options;
DROP FUNCTION IF EXISTS puzzle_game_unlock_options_update_updated_at_column();

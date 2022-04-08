-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION puzzle_games_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS puzzle_games (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    episode_id uuid NOT NULL,
    prize_pool DOUBLE PRECISION NOT NULL,
    parts_x INTEGER NOT NULL DEFAULT 5,
    parts_y INTEGER NOT NULL DEFAULT 5,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TRIGGER update_puzzle_games_modtime BEFORE
    UPDATE ON puzzle_games FOR EACH ROW EXECUTE PROCEDURE puzzle_games_update_updated_at_column();
    
-- +migrate Down
DROP TRIGGER IF EXISTS update_puzzle_games_modtime ON puzzle_games;
DROP TABLE IF EXISTS puzzle_games;
DROP FUNCTION IF EXISTS puzzle_games_update_updated_at_column();

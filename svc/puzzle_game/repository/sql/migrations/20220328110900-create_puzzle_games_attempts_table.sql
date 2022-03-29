-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION puzzle_games_attempts_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS puzzle_games_attempts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    puzzle_game_id uuid NOT NULL,
    user_id uuid NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    steps INTEGER NOT NULL DEFAULT 0,
    steps_taken INTEGER NOT NULL DEFAULT 0,
    rewards_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    result INTEGER NOT NULL DEFAULT 0,
    image VARCHAR DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TRIGGER update_puzzle_games_attempts_modtime BEFORE
    UPDATE ON puzzle_games_attempts FOR EACH ROW EXECUTE PROCEDURE puzzle_games_attempts_update_updated_at_column();
    
-- +migrate Down
DROP TRIGGER IF EXISTS update_puzzle_games_attempts_modtime ON puzzle_games_attempts;
DROP TABLE IF EXISTS puzzle_games_attempts;
DROP FUNCTION IF EXISTS puzzle_games_attempts_update_updated_at_column();

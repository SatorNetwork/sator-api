-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION puzzle_games_to_images_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- file_id (uuid == files.id)
-- puzzle_game_id (uuid == puzzle_games.id)

-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS puzzle_games_to_images (
    file_id uuid NOT NULL,
    puzzle_game_id uuid NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE TRIGGER update_puzzle_games_to_images_modtime BEFORE
    UPDATE ON puzzle_games_to_images FOR EACH ROW EXECUTE PROCEDURE puzzle_games_to_images_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_puzzle_games_to_images_modtime ON puzzle_games_to_images;
DROP TABLE IF EXISTS puzzle_games_to_images;
DROP FUNCTION IF EXISTS puzzle_games_to_images_update_updated_at_column();

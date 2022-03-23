-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION puzzle_games_to_images_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS puzzle_games_to_images (
    file_id uuid NOT NULL,
    puzzle_game_id uuid NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY(file_id) REFERENCES files(id),
    FOREIGN KEY(puzzle_game_id) REFERENCES puzzle_games(id)
    );
CREATE TRIGGER update_puzzle_games_to_images_modtime BEFORE
    UPDATE ON puzzle_games_to_images FOR EACH ROW EXECUTE PROCEDURE puzzle_games_to_images_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_puzzle_games_to_images_modtime ON puzzle_games_to_images;
DROP TABLE IF EXISTS puzzle_games_to_images;
DROP FUNCTION IF EXISTS puzzle_games_to_images_update_updated_at_column();

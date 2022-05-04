-- +migrate Up
ALTER TABLE puzzle_games_attempts DROP COLUMN result;
-- +migrate Down
ALTER TABLE puzzle_games_attempts ADD COLUMN result INTEGER NOT NULL DEFAULT 0;
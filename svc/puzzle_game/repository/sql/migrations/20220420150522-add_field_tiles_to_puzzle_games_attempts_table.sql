-- +migrate Up
ALTER TABLE
    puzzle_games_attempts
ADD
    COLUMN tiles VARCHAR;
-- +migrate Down
ALTER TABLE
    puzzle_games_attempts DROP COLUMN IF EXISTS tiles;
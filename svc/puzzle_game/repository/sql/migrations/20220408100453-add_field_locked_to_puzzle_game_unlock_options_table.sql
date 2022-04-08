-- +migrate Up
ALTER TABLE
    puzzle_game_unlock_options
ADD
    COLUMN locked BOOLEAN NOT NULL DEFAULT FALSE;

-- +migrate Down
ALTER TABLE
    puzzle_game_unlock_options DROP COLUMN IF EXISTS locked;
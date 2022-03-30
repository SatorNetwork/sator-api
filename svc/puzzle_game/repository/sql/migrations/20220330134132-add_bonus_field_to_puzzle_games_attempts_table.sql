-- +migrate Up
ALTER TABLE puzzle_games_attempts
    ADD COLUMN bonus_amount DOUBLE PRECISION NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE puzzle_games_attempts DROP COLUMN bonus_amount;
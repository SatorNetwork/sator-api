-- +migrate Up
ALTER TABLE challenges
    ADD COLUMN percent_for_quiz DOUBLE PRECISION NOT NULL DEFAULT 5,
    ADD COLUMN minimum_reward DOUBLE PRECISION NOT NULL DEFAULT 1;
-- +migrate Down
ALTER TABLE challenges
    DROP COLUMN percent_for_quiz,
    DROP COLUMN minimum_reward;

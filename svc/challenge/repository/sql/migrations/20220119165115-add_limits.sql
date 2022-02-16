-- +migrate Up
ALTER TABLE challenges
    ADD COLUMN max_winners INT,
    ADD COLUMN questions_per_game INT NOT NULL DEFAULT 5,
    ADD COLUMN min_correct_answers INT NOT NULL DEFAULT 1;
-- +migrate Down
ALTER TABLE challenges
    DROP COLUMN max_winners,
    DROP COLUMN questions_per_game,
    DROP COLUMN min_correct_answers;
-- +migrate Up
ALTER TABLE challenges
ADD COLUMN kind INT NOT NULL DEFAULT 0;
-- +migrate Down
ALTER TABLE challenges DROP COLUMN kind;
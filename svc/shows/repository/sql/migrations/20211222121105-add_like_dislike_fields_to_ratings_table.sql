-- +migrate Up
ALTER TABLE ratings
    ADD COLUMN like_dislike INT DEFAULT 0;
-- +migrate Down
ALTER TABLE ratings DROP COLUMN like_dislike;
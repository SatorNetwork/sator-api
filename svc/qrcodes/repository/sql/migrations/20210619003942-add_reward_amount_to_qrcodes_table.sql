-- +migrate Up
ALTER TABLE qrcodes
ADD COLUMN reward_amount DOUBLE PRECISION DEFAULT 0;
-- +migrate Down
ALTER TABLE qrcodes DROP COLUMN reward_amount;
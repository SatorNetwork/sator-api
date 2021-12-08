-- +migrate Up
ALTER TABLE users
    ADD COLUMN kyc_approved BOOLEAN DEFAULT FALSE;
-- +migrate Down
ALTER TABLE users DROP COLUMN kyc_approved;
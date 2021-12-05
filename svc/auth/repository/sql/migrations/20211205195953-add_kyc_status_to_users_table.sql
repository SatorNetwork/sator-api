-- +migrate Up
ALTER TABLE users
    ADD COLUMN kyc_status VARCHAR DEFAULT 'declined';
-- +migrate Down
ALTER TABLE users DROP COLUMN kyc_status;
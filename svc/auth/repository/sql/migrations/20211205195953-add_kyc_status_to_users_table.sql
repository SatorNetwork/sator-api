-- +migrate Up
ALTER TABLE users
    ADD COLUMN kyc_status VARCHAR DEFAULT 'verification_needed';
-- +migrate Down
ALTER TABLE users DROP COLUMN kyc_status;
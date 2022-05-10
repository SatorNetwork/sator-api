-- +migrate Up
ALTER TABLE users
    ADD COLUMN sanitized_email VARCHAR DEFAULT NULL,
    ADD COLUMN email_hash VARCHAR DEFAULT NULL;
    
-- +migrate Down
ALTER TABLE users
    DROP COLUMN sanitized_email,
    DROP COLUMN email_hash;
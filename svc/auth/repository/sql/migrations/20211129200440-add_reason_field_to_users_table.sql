-- +migrate Up
ALTER TABLE users
    ADD COLUMN block_reason VARCHAR DEFAULT NULL;
    
-- +migrate Down
ALTER TABLE users DROP COLUMN block_reason;
-- +migrate Up
ALTER TABLE wallets
ADD COLUMN wallet_type VARCHAR NOT NULL DEFAULT "sao";
ALTER TABLE ONLY wallets
ALTER COLUMN wallet_type
SET DEFAULT NULL;
ALTER TABLE wallets DROP COLUMN wallet_name;
-- +migrate Down
ALTER TABLE wallets
ADD COLUMN wallet_name VARCHAR DEFAULT NULL;
ALTER TABLE wallets DROP COLUMN wallet_type;
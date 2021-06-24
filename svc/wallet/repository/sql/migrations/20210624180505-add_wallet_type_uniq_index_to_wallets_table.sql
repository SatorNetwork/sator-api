-- +migrate Up
CREATE UNIQUE INDEX wallets_uniq_user_wallet_type ON wallets USING btree(user_id, wallet_type);
-- +migrate Down
DROP INDEX IF EXISTS wallets_uniq_user_wallet_type;
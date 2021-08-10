
-- +migrate Up
DROP INDEX IF EXISTS wallets_user_id;

-- +migrate Down
CREATE UNIQUE INDEX wallets_user_id ON wallets USING BTREE (user_id, solana_account_id);

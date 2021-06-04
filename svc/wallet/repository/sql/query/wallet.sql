-- name: GetWalletsByUserID :many
SELECT *
FROM wallets
WHERE user_id = $1;
-- name: CreateWallet :one
INSERT INTO wallets (user_id, solana_account_id, wallet_name)
VALUES (@user_id, @solana_account_id, @wallet_name) RETURNING *;
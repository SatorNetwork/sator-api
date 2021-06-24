-- name: GetWalletsByUserID :many
SELECT *
FROM wallets
WHERE user_id = $1
ORDER BY sort ASC;
-- name: CreateWallet :one
INSERT INTO wallets (user_id, solana_account_id, wallet_type, sort)
VALUES (
        @user_id,
        @solana_account_id,
        @wallet_type,
        @sort
    ) RETURNING *;
-- name: GetWalletBySolanaAccountID :one
SELECT *
FROM wallets
WHERE solana_account_id = $1
LIMIT 1;
-- name: GetWalletByID :one
SELECT *
FROM wallets
WHERE id = $1
LIMIT 1;
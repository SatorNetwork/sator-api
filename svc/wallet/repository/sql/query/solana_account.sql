-- name: AddSolanaAccount :one
INSERT INTO solana_accounts (
        public_key,
        private_key
    )
VALUES (
        @public_key,
        @private_key
    ) ON CONFLICT (public_key) DO NOTHING RETURNING *;
-- name: GetSolanaAccountByID :one
SELECT *
FROM solana_accounts
WHERE id = @id
LIMIT 1;
-- name: GetSolanaAccountByUserID :one
SELECT *
FROM solana_accounts
WHERE id = (
        SELECT solana_account_id
        FROM wallets
        WHERE user_id = @user_id
        LIMIT 1
    )
LIMIT 1;
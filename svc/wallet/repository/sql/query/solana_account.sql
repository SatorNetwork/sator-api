-- name: AddSolanaAccount :one
INSERT INTO solana_accounts (
        account_type,
        public_key,
        private_key
    )
VALUES (
        @account_type,
        @public_key,
        @private_key
    ) ON CONFLICT (public_key) DO NOTHING RETURNING *;
-- name: GetSolanaAccountByType :one
SELECT *
FROM solana_accounts
WHERE account_type = @account_type
LIMIT 1;
-- name: GetSolanaAccountByID :one
SELECT *
FROM solana_accounts
WHERE id = @id
LIMIT 1;
-- name: GetSolanaAccountByUserIDAndType :one
SELECT *
FROM solana_accounts
WHERE id = (
        SELECT solana_account_id
        FROM wallets
        WHERE user_id = @user_id
            AND account_type = @account_type
        LIMIT 1
    )
LIMIT 1;
-- name: GetSolanaAccountTypeByPublicKey :one
SELECT account_type
FROM solana_accounts
WHERE public_key = $1
LIMIT 1;
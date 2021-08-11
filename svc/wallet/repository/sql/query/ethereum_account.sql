-- name: AddEthereumAccount :one
INSERT INTO ethereum_accounts (
        public_key,
        private_key,
        address
    )
VALUES (
        @public_key,
        @private_key,
        @address
    ) ON CONFLICT (public_key) DO NOTHING RETURNING *;
-- name: GetEthereumAccountByID :one
SELECT *
FROM ethereum_accounts
WHERE id = @id
LIMIT 1;
-- name: GetEthereumAccountByUserIDAndType :one
SELECT *
FROM ethereum_accounts
WHERE id = (
        SELECT ethereum_account_id
        FROM wallets
        WHERE user_id = @user_id
            AND wallet_type = @wallet_type
        LIMIT 1
    )
LIMIT 1;

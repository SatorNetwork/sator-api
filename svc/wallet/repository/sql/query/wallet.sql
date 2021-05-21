-- name: GetWalletByID :one
SELECT *
FROM wallets
WHERE id = $1
LIMIT 1;
-- name: GetWalletByAssetName :one
SELECT *
FROM wallets
WHERE user_id = $1
    AND asset_name = $2
LIMIT 1;
-- name: GetWalletsByUserID :many
SELECT *
FROM wallets
WHERE user_id = $1;
-- name: CreateWallet :one
INSERT INTO wallets (
        user_id,
        asset_name,
        wallet_address,
        public_key,
        private_key,
        status
    )
VALUES (
        @user_id,
        @asset_name,
        @wallet_address,
        @public_key,
        @private_key,
        @status
    ) RETURNING *;
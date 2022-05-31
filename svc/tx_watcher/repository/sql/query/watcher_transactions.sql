-- name: RegisterTransaction :one
INSERT INTO watcher_transactions (
    serialized_message,
    latest_valid_block_height,
    account_aliases,
    tx_hash,
    status
)
VALUES (
    @serialized_message,
    @latest_valid_block_height,
    @account_aliases,
    @tx_hash,
    @status
) RETURNING *;

-- name: GetTransactionsByStatus :many
SELECT * FROM watcher_transactions
WHERE status = @status;

-- name: GetAllTransactions :many
SELECT * FROM watcher_transactions;

-- name: UpdateTransactionStatus :exec
UPDATE watcher_transactions
SET status = @status
WHERE id = @id;

-- name: UpdateTransaction :exec
UPDATE watcher_transactions
SET latest_valid_block_height = @latest_valid_block_height,
    tx_hash = @tx_hash
WHERE id = @id;

-- name: CleanTransactions :exec
DELETE FROM watcher_transactions;
-- name: AddTokenTransfer :one 
INSERT INTO token_transfers (user_id, sender_address, recipient_address, amount, tx_hash, status)
VALUES (
        @user_id,
        @sender_address,
        @recipient_address,
        @amount,
        @tx_hash,
        @status
    ) RETURNING *;

-- name: UpdateTokenTransfer :exec
UPDATE token_transfers
SET status = @status,
    tx_hash = @tx_hash
WHERE id = @id;
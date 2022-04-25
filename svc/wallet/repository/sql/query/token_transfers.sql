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

-- name: CheckRecipientAddress :one
SELECT count(DISTINCT user_id)
FROM token_transfers 
WHERE recipient_address = @recipient_address
    AND user_id != @user_id;

-- name: DoesUserHaveFraudulentTransfers :one
SELECT (count(DISTINCT user_id) > 0)::BOOLEAN as fraud_detected
FROM token_transfers 
WHERE user_id = @user_id
    AND status = 3;
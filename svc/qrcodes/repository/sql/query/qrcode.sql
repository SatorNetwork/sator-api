-- name: GetDataByQRCodeID :one
SELECT *
FROM qrcodes
WHERE id = $1
    LIMIT 1;
-- name: AddQRCode :one
INSERT INTO qrcodes (
    show_id,
    episode_id,
    starts_at,
    expires_at,
    reward_amount
)
VALUES (
           @show_id,
           @episode_id,
           @starts_at,
           @expires_at,
           @reward_amount
       ) RETURNING *;
-- name: UpdateQRCode :exec
UPDATE qrcodes
SET show_id = @show_id,
    episode_id = @episode_id,
    starts_at = @starts_at,
    expires_at = @expires_at,
    reward_amount = @reward_amount,
    updated_at = @updated_at
WHERE id = @id;
-- name: DeleteQRCodeByID :exec
DELETE FROM qrcodes
WHERE id = @id;
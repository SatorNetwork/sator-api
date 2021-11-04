-- name: GetScannedQRCodeByUserID :one
SELECT *
FROM scanned_qrcodes
WHERE user_id = $1 AND qrcode_id = $2
    LIMIT 1;
-- name: AddScannedQRCode :one
INSERT INTO scanned_qrcodes (
    user_id,
    qrcode_id
)
VALUES (
           @user_id,
           @qrcode_id
       ) RETURNING *;
-- name: DeleteScannedQRCode :exec
DELETE FROM scanned_qrcodes
WHERE user_id = $1 AND qrcode_id = $2;
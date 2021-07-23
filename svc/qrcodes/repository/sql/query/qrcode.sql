-- name: GetQRCodeDataByID :one
SELECT *
FROM qrcodes
WHERE id = $1
    LIMIT 1;
-- name: GetQRCodesData :many
SELECT *
FROM qrcodes
ORDER BY updated_at DESC,
         created_at DESC
LIMIT $1 OFFSET $2;

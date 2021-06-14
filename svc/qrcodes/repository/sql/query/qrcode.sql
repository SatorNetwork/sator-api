-- name: GetDataByQRCodeID :one
SELECT *
FROM qrcodes
WHERE id = $1
    LIMIT 1;
-- name: AddReferralCodeData :one
INSERT INTO referral_codes (
    id,
    title,
    code,
    is_personal,
    user_id
)
VALUES (
            @id,
           @title,
           @code,
           @is_personal,
           @user_id
       ) RETURNING *;
-- name: GetReferralCodeDataByUserID :many
SELECT *
FROM referral_codes
WHERE user_id = $1;
-- name: GetReferralCodesDataList :many
SELECT *
FROM referral_codes
ORDER BY created_at DESC;
-- name: UpdateReferralCodeData :exec
UPDATE referral_codes
SET title = @title,
    code = @code,
    is_personal = @is_personal,
    user_id = @user_id
WHERE id = @id;
-- name: DeleteReferralCodeDataByID :exec
DELETE FROM referral_codes
WHERE id = @id;

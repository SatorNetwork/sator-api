-- name: AddReferralCodeData :one
INSERT INTO referral_codes (
        id,
        title,
        code,
        referral_link,
        is_personal,
        user_id
    )
VALUES (
        @id,
        @title,
        @code,
        @referral_link,
        @is_personal,
        @user_id
    ) RETURNING *;
-- name: GetReferralCodeDataByUserID :one
SELECT *
FROM referral_codes
WHERE user_id = $1;
-- name: GetReferralCodeDataByCode :one
SELECT *
FROM referral_codes
WHERE code = $1 
LIMIT 1;
-- name: GetReferralCodesDataList :many
SELECT *
FROM referral_codes
ORDER BY created_at DESC
    LIMIT $1 OFFSET $2;
-- name: UpdateReferralCodeData :exec
UPDATE referral_codes
SET title = @title,
    code = @code,
    referral_link = @referral_link,
    is_personal = @is_personal,
    user_id = @user_id
WHERE id = @id;
-- name: DeleteReferralCodeDataByID :exec
DELETE FROM referral_codes
WHERE id = @id;
-- name: GetNumberOfReferralCodes :one
SELECT COUNT (id)
FROM referral_codes;

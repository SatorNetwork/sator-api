-- name: GetReferralsWithPaginationByUserID :many
SELECT *
FROM referrals
WHERE user_id = $1
ORDER BY created_at DESC
    LIMIT $2 OFFSET $3;
-- name: AddReferral :exec
INSERT INTO referrals (
    referral_code_id,
    user_id
)
VALUES (
    @referral_code_id,
    @user_id
);
-- name: GetReferralCodeByID :one
SELECT referral_code_id
FROM referrals
WHERE user_id = $1;

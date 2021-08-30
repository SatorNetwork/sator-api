-- name: GetReferralCodeByID :one
SELECT referral_code_id
FROM referrals
WHERE user_id = $1;

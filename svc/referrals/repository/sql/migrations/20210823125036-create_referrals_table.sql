-- +migrate Up
CREATE TABLE IF NOT EXISTS referrals (
    referral_code_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(referral_code_id, user_id)
    );
-- +migrate Down
DROP TABLE IF EXISTS referrals;

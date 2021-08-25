-- +migrate Up
CREATE TABLE IF NOT EXISTS referral_codes (
    id uuid NOT NULL,
    title VARCHAR,
    code VARCHAR NOT NULL,
    is_personal BOOLEAN DEFAULT FALSE,
    user_id uuid DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (id, code)
    );
-- +migrate Down
DROP TABLE IF EXISTS referral_codes;
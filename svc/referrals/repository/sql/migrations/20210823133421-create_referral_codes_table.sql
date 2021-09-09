-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS referral_codes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR,
    code VARCHAR NOT NULL,
    referral_link VARCHAR,
    is_personal BOOLEAN DEFAULT FALSE,
    user_id uuid DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX referral_codes_code ON referral_codes USING BTREE (code);
-- +migrate Down
DROP TABLE IF EXISTS referral_codes;
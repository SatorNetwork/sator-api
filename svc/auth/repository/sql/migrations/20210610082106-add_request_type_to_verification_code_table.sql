-- +migrate Up
DROP TABLE IF EXISTS user_verifications;
CREATE TABLE IF NOT EXISTS user_verifications (
    request_type INT NOT NULL,
    user_id uuid NOT NULL,
    email VARCHAR DEFAULT NULL,
    verification_code bytea NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(request_type, user_id, email)
);
CREATE INDEX email_verifications_user_id ON user_verifications USING BTREE (user_id);
CREATE INDEX email_verifications_email ON user_verifications USING BTREE (email);
-- +migrate Down
DROP TABLE IF EXISTS user_verifications;
CREATE TABLE IF NOT EXISTS user_verifications (
    user_id uuid NOT NULL,
    email VARCHAR NOT NULL,
    verification_code bytea NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(user_id, email)
);
CREATE INDEX email_verifications_user_id ON user_verifications USING BTREE (user_id);
CREATE INDEX email_verifications_email ON user_verifications USING BTREE (email);
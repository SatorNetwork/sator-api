-- +migrate Up
CREATE TABLE IF NOT EXISTS password_resets (
    user_id uuid NOT NULL,
    email VARCHAR NOT NULL,
    token bytea NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(user_id, email)
);
-- +migrate Down
DROP TABLE IF EXISTS password_resets;
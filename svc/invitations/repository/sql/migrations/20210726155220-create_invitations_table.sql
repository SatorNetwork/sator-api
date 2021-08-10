-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS invitations (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR NOT NULL,
    invited_by uuid NOT NULL,
    invited_at TIMESTAMP NOT NULL DEFAULT now(),
    accepted_by uuid DEFAULT NULL,
    accepted_at TIMESTAMP DEFAULT NULL,
    reward_received BOOLEAN
);
-- +migrate Down

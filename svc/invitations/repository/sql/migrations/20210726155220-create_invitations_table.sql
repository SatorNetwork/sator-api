-- +migrate Up
CREATE TABLE IF NOT EXISTS invitations (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    invitee_email VARCHAR NOT NULL,
    normalized_invitee_email VARCHAR NOT NULL,
    invited_at TIMESTAMP NOT NULL DEFAULT now(),
    invited_by uuid NOT NULL,
    accepted_at TIMESTAMP DEFAULT NULL,
    accepted_by uuid NOT NULL
    );
-- +migrate Down

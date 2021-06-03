-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION rewards_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS rewards (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    quiz_id uuid DEFAULT NULL,
    amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    withdrawn BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
-- TODO: user can receive reward for quiz only once
-- CREATE UNIQUE INDEX question_user_quiz ON rewards USING BTREE (user_id, quiz_id);
CREATE INDEX rewards_withdrawn ON rewards USING BTREE (withdrawn);
CREATE TRIGGER update_rewards_modtime BEFORE
UPDATE ON rewards FOR EACH ROW EXECUTE PROCEDURE rewards_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_rewards_modtime ON rewards;
DROP TABLE IF EXISTS rewards;
DROP FUNCTION IF EXISTS rewards_update_updated_at_column();
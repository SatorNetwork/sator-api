-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION challenges_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS challenges (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    show_id uuid NOT NULL,
    title VARCHAR NOT NULL,
    description VARCHAR DEFAULT NULL,
    prize_pool DECIMAL NOT NULL,
    players_to_start INT NOT NULL,
    time_per_question INT DEFAULT 10,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX challenges_show_id ON challenges USING BTREE (show_id);
CREATE INDEX challenges_created_at ON challenges USING BTREE (created_at);
CREATE TRIGGER update_challenges_modtime BEFORE
UPDATE ON challenges FOR EACH ROW EXECUTE PROCEDURE challenges_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_challenges_modtime ON challenges;
DROP TABLE IF EXISTS challenges;
DROP FUNCTION IF EXISTS challenges_update_updated_at_column();
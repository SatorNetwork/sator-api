-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION stake_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS stake (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    wallet_id uuid NOT NULL,
    stake_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    stake_duration INT DEFAULT 0,
    unstake_date TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE UNIQUE INDEX stake_user_id ON stake USING BTREE (user_id, wallet_id);
CREATE TRIGGER update_stake_modtime BEFORE
    UPDATE ON stake FOR EACH ROW EXECUTE PROCEDURE stake_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_stake_modtime ON stake;
DROP TABLE IF EXISTS stake;
DROP FUNCTION IF EXISTS stake_update_updated_at_column();
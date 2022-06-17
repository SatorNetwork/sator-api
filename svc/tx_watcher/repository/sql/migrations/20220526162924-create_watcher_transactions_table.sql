-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION watcher_transactions_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS watcher_transactions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    serialized_message VARCHAR NOT NULL,
    latest_valid_block_height BIGINT NOT NULL,
    account_aliases VARCHAR[] NOT NULL,
    tx_hash VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE TRIGGER update_watcher_transactions_modtime BEFORE
    UPDATE ON watcher_transactions FOR EACH ROW EXECUTE PROCEDURE watcher_transactions_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_watcher_transactions_modtime ON watcher_transactions;
DROP TABLE IF EXISTS watcher_transactions;
DROP FUNCTION IF EXISTS watcher_transactions_update_updated_at_column();

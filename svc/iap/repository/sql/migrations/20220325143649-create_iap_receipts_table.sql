-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION iap_receipts_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS iap_receipts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id VARCHAR NOT NULL,
    receipt_data VARCHAR NOT NULL,
    receipt_in_json VARCHAR NOT NULL,
    user_id uuid NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users(id)
    );
CREATE UNIQUE INDEX iap_receipts_transaction_id ON iap_receipts USING BTREE (transaction_id);
CREATE TRIGGER update_iap_receipts_modtime BEFORE
    UPDATE ON iap_receipts FOR EACH ROW EXECUTE PROCEDURE iap_receipts_update_updated_at_column();
-- +migrate Down
DROP INDEX IF EXISTS iap_receipts_transaction_id;
DROP TRIGGER IF EXISTS update_iap_receipts_modtime ON iap_receipts;
DROP TABLE IF EXISTS iap_receipts;
DROP FUNCTION IF EXISTS iap_receipts_update_updated_at_column();

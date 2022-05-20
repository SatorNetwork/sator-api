-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION solana_metrics_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS solana_metrics (
    provider_name VARCHAR PRIMARY KEY,
    not_available_errors INT NOT NULL,
    other_errors INT NOT NULL,
    success_calls INT NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE TRIGGER update_solana_metrics_modtime BEFORE
    UPDATE ON solana_metrics FOR EACH ROW EXECUTE PROCEDURE solana_metrics_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_solana_metrics_modtime ON solana_metrics;
DROP TABLE IF EXISTS solana_metrics;
DROP FUNCTION IF EXISTS solana_metrics_update_updated_at_column();

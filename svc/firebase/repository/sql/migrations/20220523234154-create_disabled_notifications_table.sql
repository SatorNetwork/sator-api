-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION disabled_notifications_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS disabled_notifications (
    user_id uuid NOT NULL,
    topic VARCHAR NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(user_id, topic),
    FOREIGN KEY (user_id) REFERENCES users(id)
    );
CREATE TRIGGER update_disabled_notifications_modtime BEFORE
    UPDATE ON disabled_notifications FOR EACH ROW EXECUTE PROCEDURE disabled_notifications_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_disabled_notifications_modtime ON disabled_notifications;
DROP TABLE IF EXISTS disabled_notifications;
DROP FUNCTION IF EXISTS disabled_notifications_update_updated_at_column();

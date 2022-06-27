-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION read_announcements_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS read_announcements (
    announcement_id uuid NOT NULL,
    user_id uuid NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(announcement_id, user_id),
    FOREIGN KEY (announcement_id) REFERENCES announcements(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
    );
CREATE TRIGGER update_read_announcements_modtime BEFORE
    UPDATE ON read_announcements FOR EACH ROW EXECUTE PROCEDURE read_announcements_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_read_announcements_modtime ON read_announcements;
DROP TABLE IF EXISTS read_announcements;
DROP FUNCTION IF EXISTS read_announcements_update_updated_at_column();

-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION announcements_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS announcements (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    action_url VARCHAR NOT NULL,
    starts_at TIMESTAMP NOT NULL,
    ends_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE TRIGGER update_announcements_modtime BEFORE
    UPDATE ON announcements FOR EACH ROW EXECUTE PROCEDURE announcements_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_announcements_modtime ON announcements;
DROP TABLE IF EXISTS announcements;
DROP FUNCTION IF EXISTS announcements_update_updated_at_column();

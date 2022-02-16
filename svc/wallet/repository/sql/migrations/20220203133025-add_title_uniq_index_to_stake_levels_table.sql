-- +migrate Up
CREATE UNIQUE INDEX stake_levels_uniq_title ON stake_levels USING btree(title);
-- +migrate Down
DROP INDEX IF EXISTS stake_levels_uniq_title;
-- +migrate Up
CREATE TABLE IF NOT EXISTS users_devices (
    user_id uuid NOT NULL,
    device_id VARCHAR NOT NULL,
    PRIMARY KEY(user_id, device_id)
);

-- +migrate Down
DROP TABLE IF EXISTS users_devices;
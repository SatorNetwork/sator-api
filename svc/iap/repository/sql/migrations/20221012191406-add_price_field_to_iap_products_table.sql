
-- +migrate Up
ALTER TABLE iap_products
ADD COLUMN price_in_usd DOUBLE PRECISION NOT NULL DEFAULT 0.0;

-- +migrate Down
ALTER TABLE iap_products
DROP COLUMN price_in_usd;

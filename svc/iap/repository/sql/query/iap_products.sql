-- name: CreateIapProduct :one
INSERT INTO iap_products (
    id,
    price_in_sao,
    price_in_usd
)
VALUES (
    @id,
    @price_in_sao,
    @price_in_usd
) RETURNING *;

-- name: GetIapProductByID :one
SELECT * FROM iap_products
WHERE id = @id;

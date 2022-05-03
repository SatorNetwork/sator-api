-- name: CreateIapProduct :one
INSERT INTO iap_products (
    id,
    price_in_sao
)
VALUES (
    @id,
    @price_in_sao
) RETURNING *;

-- name: GetIapProductByID :one
SELECT * FROM iap_products
WHERE id = @id;

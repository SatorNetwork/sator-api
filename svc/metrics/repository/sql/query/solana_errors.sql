-- name: RegisterProviderError :exec
INSERT INTO solana_errors (
    provider_name,
    error_message,
    counter
)
VALUES (
    @provider_name,
    @error_message,
    1
) ON CONFLICT (provider_name, error_message) DO UPDATE
SET
    counter = solana_errors.counter + 1;

-- name: GetErrorCounter :one
SELECT * FROM solana_errors
WHERE provider_name = @provider_name AND error_message = @error_message;

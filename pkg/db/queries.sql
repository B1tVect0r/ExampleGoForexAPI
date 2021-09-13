-- name: AddCurrency :one
INSERT INTO currencies(code)
VALUES($1)
ON CONFLICT(code) DO NOTHING
RETURNING *;

-- name: GetCurrencies :many
SELECT * from currencies
ORDER BY code;

-- name: GetAllExchangeRatesForCurrency :many
SELECT * from exchange_rates
WHERE from_currency=$1;

-- name: GetExchangeRatesForCurrency :many
SELECT * from exchange_rates
WHERE from_currency=$1 AND to_currency=ANY($2);

-- name: SetExchangeRate :exec
INSERT INTO exchange_rates(from_currency, to_currency, rate, rate_at)
VALUES
    ($1,$2,$3,NOW()::timestamp),
    ($2,$1,1/$3,NOW()::timestamp)
ON CONFLICT (from_currency, to_currency)
DO
    UPDATE SET rate=EXCLUDED.rate, rate_at=EXCLUDED.rate_at;

-- name: GetProjectSecret :one
SELECT hashed_secret
FROM api_keys
WHERE project_id=$1;

-- name: SetProjectSecret :one
INSERT INTO api_keys(project_id, hashed_secret)
VALUES($1,$2)
ON CONFLICT(project_id)
DO
    UPDATE SET hashed_secret=EXCLUDED.hashed_secret
returning *;
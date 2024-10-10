-- name: FindTransactionById :one
SELECT id, quantity, type, origin_location_id, destiny_location_id, material_id, created_at FROM transaction
WHERE id = $1 LIMIT 1;

-- name: FetchPaginatedTransactions :many
SELECT id, quantity, type, origin_location_id, destiny_location_id, material_id, created_at FROM transaction
ORDER BY created_at LIMIT $1 OFFSET $2;

-- name: GetTransactionTableSize :one
SELECT count(*) AS exact_count FROM transaction;

-- name: CreateTransactionWithDL :one
INSERT INTO transaction (
  quantity, type, destiny_location_id, material_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, quantity, type, origin_location_id, destiny_location_id, material_id, created_at;

-- name: CreateTransactionWithOL :one
INSERT INTO transaction (
  quantity, type, origin_location_id, material_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, quantity, type, origin_location_id, destiny_location_id, material_id, created_at;

-- name: CreateTransaction :one
INSERT INTO transaction (
  quantity, type, destiny_location_id, origin_location_id,  material_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, quantity, type, origin_location_id, destiny_location_id, material_id, created_at;

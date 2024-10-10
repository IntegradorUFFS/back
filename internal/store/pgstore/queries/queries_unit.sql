-- name: FindUnitById :one
SELECT id, name, short_name FROM unit
WHERE id = $1 LIMIT 1;

-- name: AutocompleteUnitByLikeName :many
SELECT id, name, short_name FROM unit
WHERE name ~* $1 ORDER BY name ASC LIMIT 10;

-- name: FetchPaginatedUnits :many
SELECT id, name, short_name FROM unit
ORDER BY name LIMIT $1 OFFSET $2;

-- name: GetUnitTableSize :one
SELECT count(*) AS exact_count FROM unit;

-- name: CreateUnit :one
INSERT INTO unit (name, short_name) VALUES ($1, $2)
RETURNING id, name, short_name;

-- name: UpdateUnit :one
UPDATE unit
  set name = $2,
  short_name = $3
WHERE id = $1
RETURNING id, name, short_name;

-- name: DeleteUnit :exec
DELETE FROM unit
WHERE id = $1;

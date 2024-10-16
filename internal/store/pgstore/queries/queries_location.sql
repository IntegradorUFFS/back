-- name: FindLocationById :one
SELECT id, name FROM location
WHERE id = $1 LIMIT 1;

-- name: AutocompleteLocationByLikeName :many
SELECT id, name FROM location
WHERE name ~* $1 ORDER BY name LIMIT 10;

-- name: FetchPaginatedLocations :many
SELECT id, name FROM location
ORDER BY name LIMIT $1 OFFSET $2;

-- name: GetLocationTableSize :one
SELECT count(*) AS exact_count FROM location;

-- name: CreateLocation :one
INSERT INTO location (name) VALUES ($1)
RETURNING id, name;

-- name: UpdateLocation :one
UPDATE location
  set name = $2
WHERE id = $1
RETURNING id, name;

-- name: DeleteLocation :exec
DELETE FROM location
WHERE id = $1;

-- name: FindCategoryById :one
SELECT id, name FROM category
WHERE id = $1 LIMIT 1;

-- name: AutocompleteCategoryByLikeName :many
SELECT id, name FROM category
WHERE name ~* $1 ORDER BY name LIMIT 10;

-- name: FetchPaginatedCategories :many
SELECT id, name FROM category
ORDER BY name LIMIT $1 OFFSET $2;

-- name: GetCategoryTableSize :one
SELECT count(*) AS exact_count FROM category;

-- name: CreateCategory :one
INSERT INTO category (name) VALUES ($1)
RETURNING id, name;

-- name: UpdateCategory :one
UPDATE category
  set name = $2
WHERE id = $1
RETURNING id, name;

-- name: DeleteCategory :exec
DELETE FROM category
WHERE id = $1;

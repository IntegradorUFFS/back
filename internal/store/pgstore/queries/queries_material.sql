-- name: FindMaterialById :one
SELECT id, name, description, quantity, category_id, unit_id FROM material
WHERE id = $1 LIMIT 1;

-- name: FetchMaterials :many
SELECT id, name, description, quantity, category_id, unit_id FROM material
ORDER BY name;

-- name: FetchPaginatedMaterials :many
SELECT id, name, description, quantity, category_id, unit_id FROM material
ORDER BY name LIMIT $1 OFFSET $2;

-- name: GetMaterialTableSize :one
SELECT count(*) AS exact_count FROM material;

-- name: CreateMaterial :one
INSERT INTO material (
  name, description, category_id, unit_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, description, quantity, category_id, unit_id;

-- name: UpdateMaterial :one
UPDATE material
  set name = $2,
  description = $3,
  category_id = $4,
  unit_id = $5
WHERE id = $1
RETURNING name, description, quantity, category_id, unit_id;

-- name: UpdateMaterialQuantity :exec
UPDATE material
  set quantity = $2
WHERE id = $1;

-- name: DeleteMaterial :exec
DELETE FROM material
WHERE id = $1;

-- name: FindLocationMaterialById :one
SELECT id, quantity, material_id, location_id FROM location_material
WHERE id = $1 LIMIT 1;

-- name: FindLocationMaterialByRelations :one
SELECT id, quantity, material_id, location_id FROM location_material
WHERE material_id = $1 AND location_id = $2 LIMIT 1;


-- name: FetchPaginatedLocationMaterials :many
SELECT id, quantity, material_id, location_id FROM location_material
ORDER BY id LIMIT $1 OFFSET $2;

-- name: GetLocationMaterialTableSize :one
SELECT count(*) AS exact_count FROM location_material;

-- name: CreateLocationMaterial :one
INSERT INTO location_material (
  quantity, material_id, location_id
) VALUES (
  $1, $2, $3
)
RETURNING id, quantity, material_id, location_id;

-- name: UpdateLocationMaterialLocation :one
UPDATE location_material
  set location_id = $2
WHERE id = $1
RETURNING quantity, material_id, location_id;

-- name: UpdateLocationMaterialQuantity :one
UPDATE location_material
  set quantity = $2
WHERE id = $1
RETURNING quantity, material_id, location_id;

-- name: DeleteLocationMaterial :exec
DELETE FROM location_material
WHERE id = $1;

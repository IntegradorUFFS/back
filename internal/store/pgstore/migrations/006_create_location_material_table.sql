-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS location_material (
    "id"            uuid      PRIMARY KEY   NOT NULL    DEFAULT     gen_random_uuid(),
    "quantity"      FLOAT(24) NOT NULL,
    "material_id"   UUID NOT NULL REFERENCES "material"(id),
    "location_id"   UUID NOT NULL REFERENCES "location"(id)
);
---- create above / drop below ----
DROP TABLE IF EXISTS location_material
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

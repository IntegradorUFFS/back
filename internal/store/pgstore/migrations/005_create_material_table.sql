-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS material (
    "id"            uuid      PRIMARY KEY   NOT NULL    DEFAULT     gen_random_uuid(),
    "name"          VARCHAR(255)            NOT NULL    UNIQUE,
    "description"   VARCHAR(255),
    "quantity"      FLOAT(24)               NOT NULL    DEFAULT     0,
    "category_id"   uuid                    NOT NULL REFERENCES "category"(id),
    "unit_id"       UUID                    NOT NULL REFERENCES "unit"(id)
);
---- create above / drop below ----
DROP TABLE IF EXISTS material
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

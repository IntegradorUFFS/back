-- Write your migrate up statements here
CREATE TYPE transactionType AS ENUM ('in', 'out', 'transfer');

CREATE TABLE IF NOT EXISTS transaction (
    "id"                    uuid      PRIMARY KEY   NOT NULL    DEFAULT     gen_random_uuid(),
    "quantity"              FLOAT(24)               NOT NULL    DEFAULT     0,
    created_at              timestamp               NOT NULL    DEFAULT     NOW(),
    type                    transactionType         NOT NULL,
    "origin_location_id"    UUID          REFERENCES "location"(id), --remove not null
    "destiny_location_id"   UUID          REFERENCES "location"(id), --remove not null
    "material_id"           UUID NOT NULL REFERENCES "material"(id)
);
---- create above / drop below ----
DROP TABLE IF EXISTS transaction
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS category (
    "id"            uuid      PRIMARY KEY   NOT NULL    DEFAULT     gen_random_uuid(),
    "name"          VARCHAR(255)            NOT NULL    UNIQUE
);
---- create above / drop below ----
DROP TABLE IF EXISTS category
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

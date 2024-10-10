-- Write your migrate up statements here
CREATE TYPE userRole AS ENUM ('viewer', 'manager', 'admin');

CREATE TABLE IF NOT EXISTS users(
    "id"            UUID      PRIMARY KEY   NOT NULL    DEFAULT     gen_random_uuid(),
    "email"         VARCHAR(255)            NOT NULL    UNIQUE,
    "password"      VARCHAR(255)            NOT NULL,
    "first_name"    VARCHAR(255)            NOT NULL,
    "last_name"     VARCHAR(255),
    "role"          userRole                NOT NULL    DEFAULT     'viewer',
    created_at      timestamp               NOT NULL    DEFAULT     NOW()
);

INSERT INTO users (email, password, first_name, last_name, role)
VALUES (
  'admin@gmail.com',
  '$2a$10$K0t33Uh49zvuUHrSLp/B3ed4MTQN7n05jtIpDX4qupOWeTLTaLpPe',
  'ADMIN',
  '',
  'admin'
)
ON CONFLICT (email) DO NOTHING;
---- create above / drop below ----
DROP TABLE IF EXISTS users
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

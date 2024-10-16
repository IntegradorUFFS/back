-- name: FindUserById :one
SELECT id, email, first_name, last_name, role FROM users
WHERE id = $1 LIMIT 1;

-- name: FindUserByIdWithPassword :one
SELECT id, email, first_name, last_name, role, password FROM users
WHERE id = $1 LIMIT 1;

-- name: FindUserByEmail :one
SELECT id, email, first_name, last_name, role, password FROM users
WHERE email = $1 LIMIT 1;

-- name: FetchUsers :many
SELECT id, email, first_name, last_name, role FROM users
ORDER BY first_name;

-- name: FetchPaginatedUsers :many
SELECT id, email, first_name, last_name, role FROM users
ORDER BY first_name LIMIT $1 OFFSET $2;

-- name: FetchPaginatedUsersByRole :many
SELECT id, email, first_name, last_name, role FROM users
WHERE role = $3
ORDER BY first_name LIMIT $1 OFFSET $2;

-- name: GetUserTableSize :one
SELECT count(*) AS exact_count FROM users;

-- name: GetRoledUserTableSize :one
SELECT count(*) AS exact_count FROM users WHERE role = $1;

-- name: CreateUser :one
INSERT INTO users (
  email, password, first_name, last_name
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, email, first_name, last_name, role;

-- name: CreateUserWithRole :one
INSERT INTO users (
  email, password, first_name, last_name, role
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, email, first_name, last_name, role;

-- name: UpdateUser :one
UPDATE users
  set email = $2,
  password = $3,
  first_name = $4,
  last_name = $5
WHERE id = $1
RETURNING email, first_name, last_name, role;

-- name: UpdateUserWithRole :one
UPDATE users
  set email = $2,
  password = $3,
  first_name = $4,
  last_name = $5,
  role = $6
WHERE id = $1
RETURNING email, first_name, last_name, role;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

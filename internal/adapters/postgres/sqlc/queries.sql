-- name: ListRoles :many
SELECT * FROM role;

-- name: FindRoleById :one
SELECT * FROM role WHERE id = $1;

-- name: FindRoleByName :one
SELECT * FROM role WHERE name = $1;

-- name: CreateRole :one
INSERT INTO role (name, description) VALUES ($1, $2) RETURNING *;

-- name: ListUsers :many
SELECT users.id, users.name, users.last_names, users.email, users.username, users.status, users.birth_date, users.address, users.phone, r.name as role
FROM users
INNER JOIN role r ON users.role_id = r.id
WHERE users.status = 'active';

-- name: FindUserById :one
SELECT users.id, users.name, users.last_names, users.email, users.username, users.status, users.birth_date, users.address, users.phone, r.name as role
FROM users
INNER JOIN role r ON users.role_id = r.id
WHERE users.id = $1;

-- name: FindUserByEmail :one
SELECT users.id, users.name, users.last_names, users.email, users.username, users.password, users.status, users.birth_date, users.address, users.phone, r.name as role
FROM users
INNER JOIN role r ON users.role_id = r.id
WHERE users.email = $1;

-- name: FindUserByUsername :one
SELECT users.id, users.name, users.last_names, users.email, users.username, users.password, users.status, users.birth_date, users.address, users.phone, r.name as role
FROM users
INNER JOIN role r ON users.role_id = r.id
WHERE users.username = $1;

-- name: CreateUser :one
INSERT INTO users (id, name, last_names, email, username, password, status, birth_date, address, phone, role_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;

-- name: UpdateUser :one
UPDATE users
SET name = $1, last_names = $2, username = $3
WHERE id = $4
RETURNING id;

-- name: DeleteUser :exec
UPDATE users SET status = 'deleted' WHERE id = $1;

-- name: CreateTokenBlacklist :one
INSERT INTO token_blacklist (token) VALUES ($1) RETURNING *;

-- name: FindTokenBlacklist :one
SELECT token FROM token_blacklist WHERE token = $1;

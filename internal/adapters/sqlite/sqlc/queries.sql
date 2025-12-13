-- name: ListRoles :many
SELECT * FROM role;

-- name: FindRoleById :one
SELECT * FROM role WHERE id = ?;

-- name: FindRoleByName :one
SELECT * FROM role WHERE name = ?;

-- name: CreateRole :one
INSERT INTO role (name, description) VALUES (?, ?) RETURNING *;

-- name: ListUsers :many
SELECT user.id, user.name, user.last_names, user.email, user.username, user.status, user.birth_date, user.address, user.phone, r.name as role 
FROM user 
INNER JOIN role r ON user.role_id = r.id
WHERE user.status = 'active';

-- name: FindUserById :one
SELECT user.id, user.name, user.last_names, user.email, user.username, user.status, user.birth_date, user.address, user.phone, r.name as role 
FROM user 
INNER JOIN role r ON user.role_id = r.id
WHERE user.id = ?;

-- name: FindUserByEmail :one
SELECT user.id, user.name, user.last_names, user.email, user.username, user.password, user.status, user.birth_date, user.address, user.phone, r.name as role 
FROM user 
INNER JOIN role r ON user.role_id = r.id
WHERE user.email = ?;

-- name: FindUserByUsername :one
SELECT user.id, user.name, user.last_names, user.email, user.username, user.password, user.status, user.birth_date, user.address, user.phone, r.name as role 
FROM user 
INNER JOIN role r ON user.role_id = r.id
WHERE user.username = ?;

-- name: CreateUser :one
INSERT INTO user (id, name, last_names, email, username, password, status, birth_date, address, phone, role_id) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
RETURNING id;

-- name: UpdateUser :one
UPDATE user 
SET name = ?, last_names = ?, username = ? 
WHERE id = ? 
RETURNING id;

-- name: DeleteUser :exec
UPDATE user SET status = 'deleted' WHERE id = ?;

-- name: CreateTokenBlacklist :one
INSERT INTO token_blacklist (token) VALUES (?) RETURNING *;

-- name: FindTokenBlacklist :one
SELECT token FROM token_blacklist WHERE token = ?;
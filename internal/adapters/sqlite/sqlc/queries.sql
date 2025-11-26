-- name: ListRoles :many
SELECT * FROM role;

-- name: FindRoleById :one
SELECT * FROM role WHERE id = ?;

-- name: CreateRole :one
INSERT INTO role (name, description) VALUES (?, ?) RETURNING *;

-- name: ListUsers :many
SELECT * FROM user;

-- name: FindUserById :one
SELECT * FROM user WHERE id = ?;

-- name: FindUserByEmail :one
SELECT * FROM user WHERE email = ?;

-- name: CreateUser :one
INSERT INTO user (id, name, last_names, email, username, password, status, birth_date, address, phone, role_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateUser :one
UPDATE user SET name = ?, last_names = ?, email = ?, username = ?, password = ?, status = ?, birth_date = ?, address = ?, phone = ?, role_id = ? WHERE id = ? RETURNING *;

-- name: DeleteUser :exec
UPDATE user SET status = 'deleted' WHERE id = ?;

-- name: CreateTokenBlacklist :one
INSERT INTO token_blacklist (token) VALUES (?) RETURNING *;

-- name: FindTokenBlacklist :one
SELECT * FROM token_blacklist WHERE token = ?;
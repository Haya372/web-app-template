-- name: CreateUser :exec
insert into users(id, email, password_hash, name, status_code, created_at) values ($1, $2, $3, $4, $5, $6);

-- name: FindUserByEmail :one
select id, email, password_hash, name, status_code, created_at, updated_at from users
where email = $1;

-- name: FindAllUsers :many
SELECT id, name, email, status_code, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: FindUserByID :one
SELECT id, email, password_hash, name, status_code, created_at, updated_at FROM users
WHERE id = $1;

-- name: FindPermissionsByUserID :many
SELECT p.code
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE ur.user_id = $1;

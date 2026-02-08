-- name: CreateUser :exec
insert into users(id, email, password_hash, name, status_code, created_at) values ($1, $2, $3, $4, $5, $6);

-- name: FindUserByEmail :one
select id, email, password_hash, name, status_code, created_at, updated_at from users
where email = $1;

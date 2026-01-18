-- name: CreateUser :exec
insert into users(id, email, password_hash, name, created_at) values ($1, $2, $3, $4, $5);

-- name: FindUserByEmail :one
select id, email, password_hash, name, created_at from users
where email = $1;

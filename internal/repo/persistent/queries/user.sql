-- name: ReadUser :one
SELECT id, email, name, password_hash, created_at, updated_at
FROM users;

-- name: CreateUser :exec
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3);

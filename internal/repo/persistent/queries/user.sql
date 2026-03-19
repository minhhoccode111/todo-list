-- name: ReadUser :one
SELECT id, email, username, password_hash, created_at, updated_at
FROM users;

-- name: CreateUser :exec
INSERT INTO users (email, username, password_hash)
VALUES ($1, $2, $3);

-- name: ReadUserByID :one
SELECT id, email, name, password_hash, created_at, updated_at
FROM users
WHERE id = $1;

-- name: ReadUserByEmail :one
SELECT id, email, name, password_hash, created_at, updated_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at, device_info)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ReadRefreshTokenByHash :one
SELECT id, user_id, token_hash, expires_at, device_info, created_at
FROM refresh_tokens
WHERE token_hash = $1 AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1;

-- name: RevokeAllUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();

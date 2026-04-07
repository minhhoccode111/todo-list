-- name: CreateRefreshToken :exec
insert into refresh_tokens (user_id, token_hash, expires_at, device_info)
values ($1, $2, $3, $4);

-- name: ReadRefreshToken :one
delete from refresh_tokens
where token_hash = $1 and expires_at > now()
returning user_id;

-- name: DeleteRefreshTokenById :exec
delete from refresh_tokens
where user_id = $1 and id = $2;

-- name: DeleteRefreshTokenByHash :exec
delete from refresh_tokens
where user_id = $1 and token_hash = $2;

-- name: DeleteAllRefreshTokens :exec
delete from refresh_tokens
where user_id = $1;

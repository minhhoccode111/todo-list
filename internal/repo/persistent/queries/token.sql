-- name: CreateRefreshToken :exec
insert into refresh_tokens (user_id, token_hash, expires_at, device_info)
values ($1, $2, $3, $4);

-- name: ReadRefreshToken :one
select user_id
from refresh_tokens
where token_hash = $1 and expires_at < now();

-- name: DeleteRefreshTokenById :exec
delete from refresh_tokens
where user_id = $1 and id = $2;

-- name: DeleteRefreshTokenByHash :exec
delete from refresh_tokens
where user_id = $1 and token_hash = $2;

-- name: DeleteAllRefreshTokens :exec
delete from refresh_tokens
where user_id = $1;

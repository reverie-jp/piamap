-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token_hash, expire_time)
VALUES ($1, $2, $3, sqlc.arg(expire_time)::timestamptz);

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteRefreshTokenByHash :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1 AND user_id = sqlc.arg(user_id)::ulid;

-- name: DeleteExpiredRefreshTokensByUserID :exec
DELETE FROM refresh_tokens
WHERE user_id = sqlc.arg(user_id)::ulid AND expire_time < NOW();

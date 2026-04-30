-- name: GetAuthProviderByProvider :one
SELECT * FROM user_auth_providers
WHERE provider = $1 AND provider_user_id = $2;

-- name: CreateAuthProvider :exec
INSERT INTO user_auth_providers (id, user_id, provider, provider_user_id)
VALUES ($1, $2, $3, $4);

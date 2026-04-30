-- name: GetUserByID :one
SELECT * FROM users
WHERE id = sqlc.arg(id)::ulid;

-- name: GetUserByCustomID :one
SELECT * FROM users
WHERE custom_id = $1;

-- name: ListUsersByIDs :many
SELECT * FROM users
WHERE id = ANY(@ids::text[]);

-- name: CreateUser :exec
INSERT INTO users (id, custom_id, display_name, avatar_url)
VALUES ($1, $2, $3, $4);

-- name: DeleteUser :exec
DELETE FROM users WHERE id = sqlc.arg(id)::ulid;

-- name: UpdateUserProfile :exec
UPDATE users SET
    display_name        = COALESCE(sqlc.narg(display_name), display_name),
    biography           = COALESCE(sqlc.narg(biography), biography),
    avatar_url          = COALESCE(sqlc.narg(avatar_url), avatar_url),
    hometown            = COALESCE(sqlc.narg(hometown), hometown),
    piano_started_year  = COALESCE(sqlc.narg(piano_started_year), piano_started_year),
    years_of_experience = COALESCE(sqlc.narg(years_of_experience), years_of_experience),
    update_time         = NOW()
WHERE id = sqlc.arg(id)::ulid;

-- name: UpdateUserCustomID :exec
UPDATE users SET
    custom_id             = $2,
    custom_id_change_time = NOW(),
    update_time           = NOW()
WHERE id = sqlc.arg(id)::ulid;

-- name: IsUserCurrentlyRestricted :one
SELECT EXISTS (
    SELECT 1 FROM user_restrictions
    WHERE user_id = sqlc.arg(user_id)::ulid
      AND revoked_time IS NULL
      AND suspended_until > NOW()
)::BOOLEAN AS restricted;

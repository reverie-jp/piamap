-- name: UpsertPianoPostLike :exec
INSERT INTO piano_post_likes (user_id, piano_post_id)
VALUES (sqlc.arg(user_id)::ulid, sqlc.arg(piano_post_id)::ulid)
ON CONFLICT DO NOTHING;

-- name: DeletePianoPostLike :exec
DELETE FROM piano_post_likes
WHERE user_id = sqlc.arg(user_id)::ulid
  AND piano_post_id = sqlc.arg(piano_post_id)::ulid;

-- name: ListLikedPianoPostIDsByUser :many
-- 指定ユーザーがいいねした投稿の ID 一覧 (新しい順)。
-- create_time keyset。MVP は piano_post_id (ULID) を keyset 兼テキスト cursor として使う。
SELECT piano_post_id, create_time
FROM piano_post_likes
WHERE user_id = sqlc.arg(user_id)::ulid
  AND (sqlc.narg(after_post_id)::ulid IS NULL OR piano_post_id < sqlc.narg(after_post_id)::ulid)
ORDER BY create_time DESC, piano_post_id DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: ListLikedPostIDsForUserAndPosts :many
-- 指定ユーザーが [post_ids] のうちどの投稿にいいねしているかを返す (hydrate 用)。
-- pgx に ulid[] のエンコーダーが無いため、引数は text[] で受けて SQL 内で ulid[] にキャスト。
SELECT piano_post_id
FROM piano_post_likes
WHERE user_id = sqlc.arg(user_id)::ulid
  AND piano_post_id = ANY(sqlc.arg(piano_post_ids)::text[]::ulid[]);

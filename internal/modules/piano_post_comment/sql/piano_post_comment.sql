-- name: InsertPianoPostComment :exec
INSERT INTO piano_post_comments (id, piano_post_id, user_id, parent_comment_id, body)
VALUES (
    sqlc.arg(id)::ulid,
    sqlc.arg(piano_post_id)::ulid,
    sqlc.arg(user_id)::ulid,
    sqlc.narg(parent_comment_id)::ulid,
    sqlc.arg(body)::text
);

-- name: GetPianoPostComment :one
SELECT id, piano_post_id, user_id, parent_comment_id, body, create_time, update_time
FROM piano_post_comments
WHERE id = sqlc.arg(id)::ulid;

-- name: ListPianoPostCommentsByPost :many
-- 投稿に対するコメント一覧 (古い順)。MVP は ID (ULID) を keyset として使う。
SELECT id, piano_post_id, user_id, parent_comment_id, body, create_time, update_time
FROM piano_post_comments
WHERE piano_post_id = sqlc.arg(piano_post_id)::ulid
  AND (sqlc.narg(after_id)::ulid IS NULL OR id > sqlc.narg(after_id)::ulid)
ORDER BY id ASC
LIMIT sqlc.arg(limit_count)::int;

-- name: ListPianoPostCommentsByUser :many
-- 指定ユーザーが書いたコメント一覧 (新しい順)。
SELECT id, piano_post_id, user_id, parent_comment_id, body, create_time, update_time
FROM piano_post_comments
WHERE user_id = sqlc.arg(user_id)::ulid
  AND (sqlc.narg(after_id)::ulid IS NULL OR id < sqlc.narg(after_id)::ulid)
ORDER BY id DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: DeletePianoPostComment :exec
DELETE FROM piano_post_comments
WHERE id = sqlc.arg(id)::ulid;

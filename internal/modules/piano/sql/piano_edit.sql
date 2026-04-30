-- name: CreatePianoEdit :exec
INSERT INTO piano_edits (id, piano_id, editor_user_id, operation, changes, summary)
VALUES (
    sqlc.arg(id)::ulid,
    sqlc.arg(piano_id)::ulid,
    sqlc.narg(editor_user_id)::ulid,
    sqlc.arg(operation)::piano_edit_operation,
    sqlc.narg(changes),
    sqlc.narg(summary)
);

-- name: ListPianoEditsByPiano :many
-- 指定ピアノの編集ログを新しい順に。AIP-158: id keyset で次ページ。
SELECT
    id,
    piano_id,
    editor_user_id,
    operation,
    changes,
    summary,
    create_time
FROM piano_edits
WHERE piano_id = sqlc.arg(piano_id)::ulid
  AND (sqlc.narg(after_id)::ulid IS NULL OR id < sqlc.narg(after_id)::ulid)
ORDER BY id DESC
LIMIT sqlc.arg(limit_count)::int;

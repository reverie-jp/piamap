-- name: UpsertPianoUserList :exec
INSERT INTO piano_user_lists (user_id, piano_id, list_kind)
VALUES (sqlc.arg(user_id)::ulid, sqlc.arg(piano_id)::ulid, sqlc.arg(list_kind)::piano_list_kind)
ON CONFLICT DO NOTHING;

-- name: DeletePianoUserList :exec
DELETE FROM piano_user_lists
WHERE user_id = sqlc.arg(user_id)::ulid
  AND piano_id = sqlc.arg(piano_id)::ulid
  AND list_kind = sqlc.arg(list_kind)::piano_list_kind;

-- name: ListPianoUserListsByUser :many
-- 指定ユーザーの指定 list_kind の piano_id 一覧。create_time 単独で keyset (同点は ULID で安定化)。
-- page_token は base64 不要、最後の create_time を ISO 文字列で素直に渡す形で十分。MVP は piano_id keyset で済ませる。
SELECT piano_id, create_time
FROM piano_user_lists
WHERE user_id = sqlc.arg(user_id)::ulid
  AND list_kind = sqlc.arg(list_kind)::piano_list_kind
  AND (sqlc.narg(after_piano_id)::ulid IS NULL OR piano_id < sqlc.narg(after_piano_id)::ulid)
ORDER BY create_time DESC, piano_id DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: ListMyListKindsForPiano :many
-- 認証済みユーザーが指定ピアノに付けているリスト種別。
SELECT list_kind
FROM piano_user_lists
WHERE user_id = sqlc.arg(user_id)::ulid
  AND piano_id = sqlc.arg(piano_id)::ulid;

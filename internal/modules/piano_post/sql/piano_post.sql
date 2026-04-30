-- name: GetPianoPostByID :one
SELECT
    id,
    user_id,
    piano_id,
    visit_time,
    rating,
    body,
    ambient_noise,
    foot_traffic,
    resonance,
    key_touch_weight,
    tuning_quality,
    visibility,
    comment_count,
    like_count,
    create_time,
    update_time
FROM piano_posts
WHERE id = sqlc.arg(id)::ulid;

-- name: ListPianoPostsByPiano :many
-- 指定ピアノの投稿。public のみ (本人投稿は別途 user 経由で見る運用)。
-- AIP-158: page_token は最後の (create_time, id) を opaque に。MVP は id (ULID) を keyset に使う。
SELECT
    id,
    user_id,
    piano_id,
    visit_time,
    rating,
    body,
    ambient_noise,
    foot_traffic,
    resonance,
    key_touch_weight,
    tuning_quality,
    visibility,
    comment_count,
    like_count,
    create_time,
    update_time
FROM piano_posts
WHERE piano_id = sqlc.arg(piano_id)::ulid
  AND visibility = 'public'
  AND (sqlc.narg(after_id)::ulid IS NULL OR id < sqlc.narg(after_id)::ulid)
ORDER BY id DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: ListPianoPostsByUser :many
-- 指定ユーザーの投稿。include_private = true で private も含める (本人閲覧時)。
SELECT
    id,
    user_id,
    piano_id,
    visit_time,
    rating,
    body,
    ambient_noise,
    foot_traffic,
    resonance,
    key_touch_weight,
    tuning_quality,
    visibility,
    comment_count,
    like_count,
    create_time,
    update_time
FROM piano_posts
WHERE user_id = sqlc.arg(user_id)::ulid
  AND (sqlc.arg(include_private)::bool OR visibility = 'public')
  AND (sqlc.narg(after_id)::ulid IS NULL OR id < sqlc.narg(after_id)::ulid)
ORDER BY id DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: ListPublicPianoPosts :many
-- グローバルなタイムライン用。最新の public 投稿。
SELECT
    id,
    user_id,
    piano_id,
    visit_time,
    rating,
    body,
    ambient_noise,
    foot_traffic,
    resonance,
    key_touch_weight,
    tuning_quality,
    visibility,
    comment_count,
    like_count,
    create_time,
    update_time
FROM piano_posts
WHERE visibility = 'public'
  AND (sqlc.narg(after_id)::ulid IS NULL OR id < sqlc.narg(after_id)::ulid)
ORDER BY id DESC
LIMIT sqlc.arg(limit_count)::int;

-- name: CreatePianoPost :exec
INSERT INTO piano_posts (
    id,
    user_id,
    piano_id,
    visit_time,
    rating,
    body,
    ambient_noise,
    foot_traffic,
    resonance,
    key_touch_weight,
    tuning_quality,
    visibility
) VALUES (
    sqlc.arg(id)::ulid,
    sqlc.arg(user_id)::ulid,
    sqlc.arg(piano_id)::ulid,
    sqlc.arg(visit_time)::timestamptz,
    sqlc.arg(rating)::smallint,
    sqlc.narg(body),
    sqlc.narg(ambient_noise)::smallint,
    sqlc.narg(foot_traffic)::smallint,
    sqlc.narg(resonance)::smallint,
    sqlc.narg(key_touch_weight)::smallint,
    sqlc.narg(tuning_quality)::smallint,
    sqlc.arg(visibility)::post_visibility
);

-- name: UpdatePianoPost :exec
-- set_X が true のフィールドだけ更新する (NULL 化も含めて値そのまま反映)。
-- false のフィールドは既存値を保持。rating / visit_time / visibility は NOT NULL なので
-- usecase 層で set_X=true のとき値が必ず存在することを保証する。
UPDATE piano_posts SET
    visit_time       = CASE WHEN sqlc.arg(set_visit_time)::bool       THEN sqlc.narg(visit_time)::timestamptz       ELSE visit_time       END,
    rating           = CASE WHEN sqlc.arg(set_rating)::bool           THEN sqlc.narg(rating)::smallint              ELSE rating           END,
    body             = CASE WHEN sqlc.arg(set_body)::bool             THEN sqlc.narg(body)                          ELSE body             END,
    ambient_noise    = CASE WHEN sqlc.arg(set_ambient_noise)::bool    THEN sqlc.narg(ambient_noise)::smallint       ELSE ambient_noise    END,
    foot_traffic     = CASE WHEN sqlc.arg(set_foot_traffic)::bool     THEN sqlc.narg(foot_traffic)::smallint        ELSE foot_traffic     END,
    resonance        = CASE WHEN sqlc.arg(set_resonance)::bool        THEN sqlc.narg(resonance)::smallint           ELSE resonance        END,
    key_touch_weight = CASE WHEN sqlc.arg(set_key_touch_weight)::bool THEN sqlc.narg(key_touch_weight)::smallint    ELSE key_touch_weight END,
    tuning_quality   = CASE WHEN sqlc.arg(set_tuning_quality)::bool   THEN sqlc.narg(tuning_quality)::smallint      ELSE tuning_quality   END,
    visibility       = CASE WHEN sqlc.arg(set_visibility)::bool       THEN sqlc.narg(visibility)::post_visibility   ELSE visibility       END,
    update_time      = NOW()
WHERE id = sqlc.arg(id)::ulid;

-- name: DeletePianoPost :exec
DELETE FROM piano_posts WHERE id = sqlc.arg(id)::ulid;

-- name: UpsertPianoUserListVisited :exec
-- piano_post 作成時に「行ったことある」リストに UPSERT。冪等。
INSERT INTO piano_user_lists (user_id, piano_id, list_kind)
VALUES (sqlc.arg(user_id)::ulid, sqlc.arg(piano_id)::ulid, 'visited')
ON CONFLICT DO NOTHING;

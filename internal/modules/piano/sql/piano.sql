-- name: GetPianoByID :one
SELECT
    id,
    name,
    description,
    ST_Y(location::geometry)::float8 AS latitude,
    ST_X(location::geometry)::float8 AS longitude,
    address,
    prefecture,
    city,
    kind,
    venue_type,
    piano_type,
    piano_brand,
    piano_model,
    manufacture_year,
    hours,
    status,
    availability,
    availability_note,
    install_time,
    remove_time,
    creator_user_id,
    post_count,
    rating_sum,
    ambient_noise_count,
    ambient_noise_sum,
    foot_traffic_count,
    foot_traffic_sum,
    resonance_count,
    resonance_sum,
    key_touch_weight_count,
    key_touch_weight_sum,
    tuning_quality_count,
    tuning_quality_sum,
    wishlist_count,
    visited_count,
    favorite_count,
    create_time,
    update_time,
    0::float8 AS distance_m
FROM pianos
WHERE id = sqlc.arg(id)::ulid;

-- name: ListPianosInBBox :many
-- 表示エリア内のアクティブなピアノを返す。GIST(location) の && で bbox 包含判定。
SELECT
    id,
    name,
    description,
    ST_Y(location::geometry)::float8 AS latitude,
    ST_X(location::geometry)::float8 AS longitude,
    address,
    prefecture,
    city,
    kind,
    venue_type,
    piano_type,
    piano_brand,
    piano_model,
    manufacture_year,
    hours,
    status,
    availability,
    availability_note,
    install_time,
    remove_time,
    creator_user_id,
    post_count,
    rating_sum,
    ambient_noise_count,
    ambient_noise_sum,
    foot_traffic_count,
    foot_traffic_sum,
    resonance_count,
    resonance_sum,
    key_touch_weight_count,
    key_touch_weight_sum,
    tuning_quality_count,
    tuning_quality_sum,
    wishlist_count,
    visited_count,
    favorite_count,
    create_time,
    update_time,
    0::float8 AS distance_m
FROM pianos
WHERE status = 'active'
  AND location && ST_MakeEnvelope(
      sqlc.arg(min_lng)::float8,
      sqlc.arg(min_lat)::float8,
      sqlc.arg(max_lng)::float8,
      sqlc.arg(max_lat)::float8,
      4326
  )::geography
  AND (sqlc.narg(kind)::piano_kind IS NULL OR kind = sqlc.narg(kind)::piano_kind)
  AND (sqlc.narg(piano_type)::piano_type IS NULL OR piano_type = sqlc.narg(piano_type)::piano_type)
  AND (sqlc.narg(piano_brand)::text IS NULL OR piano_brand ILIKE sqlc.narg(piano_brand)::text)
  AND (
      sqlc.narg(min_rating_average)::float8 IS NULL
      OR (post_count > 0 AND rating_sum::float8 / post_count >= sqlc.narg(min_rating_average)::float8)
  )
  AND (
      sqlc.narg(min_ambient_noise_average)::float8 IS NULL
      OR (ambient_noise_count > 0 AND ambient_noise_sum::float8 / ambient_noise_count >= sqlc.narg(min_ambient_noise_average)::float8)
  )
  AND (
      sqlc.narg(min_foot_traffic_average)::float8 IS NULL
      OR (foot_traffic_count > 0 AND foot_traffic_sum::float8 / foot_traffic_count >= sqlc.narg(min_foot_traffic_average)::float8)
  )
  AND (
      sqlc.narg(min_resonance_average)::float8 IS NULL
      OR (resonance_count > 0 AND resonance_sum::float8 / resonance_count >= sqlc.narg(min_resonance_average)::float8)
  )
  AND (
      sqlc.narg(min_key_touch_weight_average)::float8 IS NULL
      OR (key_touch_weight_count > 0 AND key_touch_weight_sum::float8 / key_touch_weight_count >= sqlc.narg(min_key_touch_weight_average)::float8)
  )
  AND (
      sqlc.narg(min_tuning_quality_average)::float8 IS NULL
      OR (tuning_quality_count > 0 AND tuning_quality_sum::float8 / tuning_quality_count >= sqlc.narg(min_tuning_quality_average)::float8)
  )
ORDER BY post_count DESC, id
LIMIT sqlc.arg(limit_count)::int;

-- name: SearchPianosByText :many
-- ピアノ名 (name) に対する部分一致検索 (グローバル)。
-- bounds なしで使う想定。MVP は ILIKE で十分 (件数が少ない、PG_TRGM 等は将来)。
SELECT
    id,
    name,
    description,
    ST_Y(location::geometry)::float8 AS latitude,
    ST_X(location::geometry)::float8 AS longitude,
    address,
    prefecture,
    city,
    kind,
    venue_type,
    piano_type,
    piano_brand,
    piano_model,
    manufacture_year,
    hours,
    status,
    availability,
    availability_note,
    install_time,
    remove_time,
    creator_user_id,
    post_count,
    rating_sum,
    ambient_noise_count,
    ambient_noise_sum,
    foot_traffic_count,
    foot_traffic_sum,
    resonance_count,
    resonance_sum,
    key_touch_weight_count,
    key_touch_weight_sum,
    tuning_quality_count,
    tuning_quality_sum,
    wishlist_count,
    visited_count,
    favorite_count,
    create_time,
    update_time,
    0::float8 AS distance_m
FROM pianos
WHERE status = 'active'
  AND name ILIKE sqlc.arg(query_pattern)::text
  AND (sqlc.narg(kind)::piano_kind IS NULL OR kind = sqlc.narg(kind)::piano_kind)
  AND (sqlc.narg(piano_type)::piano_type IS NULL OR piano_type = sqlc.narg(piano_type)::piano_type)
  AND (sqlc.narg(piano_brand)::text IS NULL OR piano_brand ILIKE sqlc.narg(piano_brand)::text)
  AND (
      sqlc.narg(min_rating_average)::float8 IS NULL
      OR (post_count > 0 AND rating_sum::float8 / post_count >= sqlc.narg(min_rating_average)::float8)
  )
  AND (
      sqlc.narg(min_ambient_noise_average)::float8 IS NULL
      OR (ambient_noise_count > 0 AND ambient_noise_sum::float8 / ambient_noise_count >= sqlc.narg(min_ambient_noise_average)::float8)
  )
  AND (
      sqlc.narg(min_foot_traffic_average)::float8 IS NULL
      OR (foot_traffic_count > 0 AND foot_traffic_sum::float8 / foot_traffic_count >= sqlc.narg(min_foot_traffic_average)::float8)
  )
  AND (
      sqlc.narg(min_resonance_average)::float8 IS NULL
      OR (resonance_count > 0 AND resonance_sum::float8 / resonance_count >= sqlc.narg(min_resonance_average)::float8)
  )
  AND (
      sqlc.narg(min_key_touch_weight_average)::float8 IS NULL
      OR (key_touch_weight_count > 0 AND key_touch_weight_sum::float8 / key_touch_weight_count >= sqlc.narg(min_key_touch_weight_average)::float8)
  )
  AND (
      sqlc.narg(min_tuning_quality_average)::float8 IS NULL
      OR (tuning_quality_count > 0 AND tuning_quality_sum::float8 / tuning_quality_count >= sqlc.narg(min_tuning_quality_average)::float8)
  )
ORDER BY post_count DESC, id
LIMIT sqlc.arg(limit_count)::int;

-- name: ListPianosNearby :many
-- 中心点から半径 radius_m 以内のアクティブなピアノを距離順に返す。
SELECT
    id,
    name,
    description,
    ST_Y(location::geometry)::float8 AS latitude,
    ST_X(location::geometry)::float8 AS longitude,
    address,
    prefecture,
    city,
    kind,
    venue_type,
    piano_type,
    piano_brand,
    piano_model,
    manufacture_year,
    hours,
    status,
    availability,
    availability_note,
    install_time,
    remove_time,
    creator_user_id,
    post_count,
    rating_sum,
    ambient_noise_count,
    ambient_noise_sum,
    foot_traffic_count,
    foot_traffic_sum,
    resonance_count,
    resonance_sum,
    key_touch_weight_count,
    key_touch_weight_sum,
    tuning_quality_count,
    tuning_quality_sum,
    wishlist_count,
    visited_count,
    favorite_count,
    create_time,
    update_time,
    ST_Distance(
        location,
        ST_SetSRID(
            ST_MakePoint(sqlc.arg(center_lng)::float8, sqlc.arg(center_lat)::float8),
            4326
        )::geography
    )::float8 AS distance_m
FROM pianos
WHERE status = 'active'
  AND ST_DWithin(
      location,
      ST_SetSRID(
          ST_MakePoint(sqlc.arg(center_lng)::float8, sqlc.arg(center_lat)::float8),
          4326
      )::geography,
      sqlc.arg(radius_m)::float8
  )
  AND (sqlc.narg(kind)::piano_kind IS NULL OR kind = sqlc.narg(kind)::piano_kind)
  AND (sqlc.narg(piano_type)::piano_type IS NULL OR piano_type = sqlc.narg(piano_type)::piano_type)
  AND (sqlc.narg(piano_brand)::text IS NULL OR piano_brand ILIKE sqlc.narg(piano_brand)::text)
  AND (
      sqlc.narg(min_rating_average)::float8 IS NULL
      OR (post_count > 0 AND rating_sum::float8 / post_count >= sqlc.narg(min_rating_average)::float8)
  )
  AND (
      sqlc.narg(min_ambient_noise_average)::float8 IS NULL
      OR (ambient_noise_count > 0 AND ambient_noise_sum::float8 / ambient_noise_count >= sqlc.narg(min_ambient_noise_average)::float8)
  )
  AND (
      sqlc.narg(min_foot_traffic_average)::float8 IS NULL
      OR (foot_traffic_count > 0 AND foot_traffic_sum::float8 / foot_traffic_count >= sqlc.narg(min_foot_traffic_average)::float8)
  )
  AND (
      sqlc.narg(min_resonance_average)::float8 IS NULL
      OR (resonance_count > 0 AND resonance_sum::float8 / resonance_count >= sqlc.narg(min_resonance_average)::float8)
  )
  AND (
      sqlc.narg(min_key_touch_weight_average)::float8 IS NULL
      OR (key_touch_weight_count > 0 AND key_touch_weight_sum::float8 / key_touch_weight_count >= sqlc.narg(min_key_touch_weight_average)::float8)
  )
  AND (
      sqlc.narg(min_tuning_quality_average)::float8 IS NULL
      OR (tuning_quality_count > 0 AND tuning_quality_sum::float8 / tuning_quality_count >= sqlc.narg(min_tuning_quality_average)::float8)
  )
ORDER BY distance_m ASC
LIMIT sqlc.arg(limit_count)::int;

-- name: CreatePiano :exec
INSERT INTO pianos (
    id,
    name,
    description,
    location,
    address,
    prefecture,
    city,
    kind,
    venue_type,
    piano_type,
    piano_brand,
    piano_model,
    manufacture_year,
    hours,
    availability,
    availability_note,
    install_time,
    creator_user_id
) VALUES (
    sqlc.arg(id)::ulid,
    sqlc.arg(name),
    sqlc.narg(description),
    ST_SetSRID(
        ST_MakePoint(sqlc.arg(longitude)::float8, sqlc.arg(latitude)::float8),
        4326
    )::geography,
    sqlc.narg(address),
    sqlc.narg(prefecture),
    sqlc.narg(city),
    sqlc.arg(kind)::piano_kind,
    sqlc.narg(venue_type),
    sqlc.arg(piano_type)::piano_type,
    sqlc.arg(piano_brand),
    sqlc.narg(piano_model),
    sqlc.narg(manufacture_year),
    sqlc.narg(hours),
    sqlc.arg(availability)::piano_availability,
    sqlc.narg(availability_note),
    sqlc.narg(install_time)::timestamptz,
    sqlc.narg(creator_user_id)::ulid
);

-- name: UpdatePiano :exec
-- set_X が true のフィールドだけ更新する (NULL 化も含めて値そのまま反映)。
-- false のフィールドは既存値を保持。NOT NULL カラム (name / kind / piano_type / piano_brand /
-- status / availability) は usecase 層で set_X=true のとき値存在を保証する。
-- location は geography 関数のため別 :exec で更新する。
UPDATE pianos SET
    name              = CASE WHEN sqlc.arg(set_name)::bool              THEN sqlc.narg(name)                              ELSE name              END,
    description       = CASE WHEN sqlc.arg(set_description)::bool       THEN sqlc.narg(description)                       ELSE description       END,
    address           = CASE WHEN sqlc.arg(set_address)::bool           THEN sqlc.narg(address)                           ELSE address           END,
    prefecture        = CASE WHEN sqlc.arg(set_prefecture)::bool        THEN sqlc.narg(prefecture)                        ELSE prefecture        END,
    city              = CASE WHEN sqlc.arg(set_city)::bool              THEN sqlc.narg(city)                              ELSE city              END,
    kind              = CASE WHEN sqlc.arg(set_kind)::bool              THEN sqlc.narg(kind)::piano_kind                  ELSE kind              END,
    venue_type        = CASE WHEN sqlc.arg(set_venue_type)::bool        THEN sqlc.narg(venue_type)                        ELSE venue_type        END,
    piano_type        = CASE WHEN sqlc.arg(set_piano_type)::bool        THEN sqlc.narg(piano_type)::piano_type            ELSE piano_type        END,
    piano_brand       = CASE WHEN sqlc.arg(set_piano_brand)::bool       THEN sqlc.narg(piano_brand)                       ELSE piano_brand       END,
    piano_model       = CASE WHEN sqlc.arg(set_piano_model)::bool       THEN sqlc.narg(piano_model)                       ELSE piano_model       END,
    manufacture_year  = CASE WHEN sqlc.arg(set_manufacture_year)::bool  THEN sqlc.narg(manufacture_year)                  ELSE manufacture_year  END,
    hours             = CASE WHEN sqlc.arg(set_hours)::bool             THEN sqlc.narg(hours)                             ELSE hours             END,
    status            = CASE WHEN sqlc.arg(set_status)::bool            THEN sqlc.narg(status)::piano_status              ELSE status            END,
    availability      = CASE WHEN sqlc.arg(set_availability)::bool      THEN sqlc.narg(availability)::piano_availability  ELSE availability      END,
    availability_note = CASE WHEN sqlc.arg(set_availability_note)::bool THEN sqlc.narg(availability_note)                 ELSE availability_note END,
    install_time      = CASE WHEN sqlc.arg(set_install_time)::bool      THEN sqlc.narg(install_time)::timestamptz         ELSE install_time      END,
    remove_time       = CASE WHEN sqlc.arg(set_remove_time)::bool       THEN sqlc.narg(remove_time)::timestamptz          ELSE remove_time       END,
    update_time       = NOW()
WHERE id = sqlc.arg(id)::ulid;

-- name: UpdatePianoLocation :exec
UPDATE pianos SET
    location = ST_SetSRID(
        ST_MakePoint(sqlc.arg(longitude)::float8, sqlc.arg(latitude)::float8),
        4326
    )::geography,
    update_time = NOW()
WHERE id = sqlc.arg(id)::ulid;

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
  AND (
      sqlc.narg(min_rating_average)::float8 IS NULL
      OR (post_count > 0 AND rating_sum::float8 / post_count >= sqlc.narg(min_rating_average)::float8)
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
  AND (
      sqlc.narg(min_rating_average)::float8 IS NULL
      OR (post_count > 0 AND rating_sum::float8 / post_count >= sqlc.narg(min_rating_average)::float8)
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
-- COALESCE で「指定されたフィールドだけ更新」を実現。location は別 :exec で更新する (geography 関数のため)。
UPDATE pianos SET
    name              = COALESCE(sqlc.narg(name), name),
    description       = COALESCE(sqlc.narg(description), description),
    address           = COALESCE(sqlc.narg(address), address),
    prefecture        = COALESCE(sqlc.narg(prefecture), prefecture),
    city              = COALESCE(sqlc.narg(city), city),
    kind              = COALESCE(sqlc.narg(kind)::piano_kind, kind),
    venue_type        = COALESCE(sqlc.narg(venue_type), venue_type),
    piano_type        = COALESCE(sqlc.narg(piano_type)::piano_type, piano_type),
    piano_brand       = COALESCE(sqlc.narg(piano_brand), piano_brand),
    piano_model       = COALESCE(sqlc.narg(piano_model), piano_model),
    manufacture_year  = COALESCE(sqlc.narg(manufacture_year), manufacture_year),
    hours             = COALESCE(sqlc.narg(hours), hours),
    status            = COALESCE(sqlc.narg(status)::piano_status, status),
    availability      = COALESCE(sqlc.narg(availability)::piano_availability, availability),
    availability_note = COALESCE(sqlc.narg(availability_note), availability_note),
    install_time      = COALESCE(sqlc.narg(install_time)::timestamptz, install_time),
    remove_time       = COALESCE(sqlc.narg(remove_time)::timestamptz, remove_time),
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

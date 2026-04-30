CREATE EXTENSION IF NOT EXISTS postgis;

CREATE DOMAIN ulid AS TEXT CHECK (LENGTH(VALUE) = 26);

-- ============================================================================
-- Auth & Users
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id ulid PRIMARY KEY,
    -- @username 相当。URL 識別子・メンション対象。a-z / 0-9 / _ の 3-20 文字
    custom_id VARCHAR(20) NOT NULL UNIQUE CHECK (custom_id ~ '^[a-z0-9_]{3,20}$'),
    custom_id_change_time TIMESTAMPTZ,
    display_name VARCHAR(30) NOT NULL DEFAULT 'unknown',
    biography TEXT,
    avatar_url TEXT,
    hometown VARCHAR(80),
    piano_started_year SMALLINT,
    years_of_experience SMALLINT,
    -- 信頼ライン判定用の集計列(piano_posts / piano_edits / reports トリガで保守)
    post_count INT NOT NULL DEFAULT 0,
    edit_count INT NOT NULL DEFAULT 0,
    -- 自分が target_type='user' として通報された件数のうち status='resolved' のもの
    report_received_count INT NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE auth_provider AS ENUM ('google');

CREATE TABLE IF NOT EXISTS user_auth_providers (
    id ulid PRIMARY KEY,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider auth_provider NOT NULL,
    provider_user_id TEXT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_user_auth_providers_user_id ON user_auth_providers(user_id);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id ulid PRIMARY KEY,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expire_time TIMESTAMPTZ NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- ============================================================================
-- Pianos
-- ============================================================================

CREATE TYPE piano_status AS ENUM ('pending', 'active', 'temporary', 'removed');
CREATE TYPE piano_type AS ENUM ('grand', 'upright', 'electronic', 'unknown');
CREATE TYPE piano_kind AS ENUM ('street', 'practice_room', 'other');
CREATE TYPE piano_list_kind AS ENUM ('wishlist', 'visited', 'favorite');
CREATE TYPE piano_availability AS ENUM (
    'regular',
    'irregular',
    'event_only',
    'weather_dependent'
);

CREATE TYPE piano_edit_operation AS ENUM (
    'create',
    'update',
    'photo_add',
    'photo_remove',
    'status_change',
    'kind_change',
    'restore'
);

CREATE TABLE IF NOT EXISTS pianos (
    id ulid PRIMARY KEY,
    name VARCHAR(80) NOT NULL,
    description TEXT,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    address TEXT,
    prefecture VARCHAR(20),
    city VARCHAR(40),
    kind piano_kind NOT NULL DEFAULT 'street',
    venue_type VARCHAR(50),
    piano_type piano_type NOT NULL DEFAULT 'unknown',
    piano_brand VARCHAR(50) NOT NULL DEFAULT 'unknown',
    piano_model VARCHAR(50),
    manufacture_year SMALLINT,
    hours TEXT,
    status piano_status NOT NULL DEFAULT 'active',
    availability piano_availability NOT NULL DEFAULT 'regular',
    availability_note TEXT,
    install_time TIMESTAMPTZ,
    remove_time TIMESTAMPTZ,
    creator_user_id ulid REFERENCES users(id) ON DELETE SET NULL,
    -- 集計列(piano_posts トリガで保守)
    --   post_count: 総投稿数 = rating付き投稿数(rating は必須なので両者は同一)
    --   平均評価は rating_sum / NULLIF(post_count, 0)
    --   5属性は任意なので null-aware で個別に count/sum を保守
    post_count INT NOT NULL DEFAULT 0,
    rating_sum INT NOT NULL DEFAULT 0,
    ambient_noise_count INT NOT NULL DEFAULT 0,
    ambient_noise_sum INT NOT NULL DEFAULT 0,
    foot_traffic_count INT NOT NULL DEFAULT 0,
    foot_traffic_sum INT NOT NULL DEFAULT 0,
    resonance_count INT NOT NULL DEFAULT 0,
    resonance_sum INT NOT NULL DEFAULT 0,
    key_touch_weight_count INT NOT NULL DEFAULT 0,
    key_touch_weight_sum INT NOT NULL DEFAULT 0,
    tuning_quality_count INT NOT NULL DEFAULT 0,
    tuning_quality_sum INT NOT NULL DEFAULT 0,
    -- ユーザーリスト集計列(piano_user_lists トリガで保守)
    wishlist_count INT NOT NULL DEFAULT 0,
    visited_count INT NOT NULL DEFAULT 0,
    favorite_count INT NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pianos_location ON pianos USING GIST (location);
CREATE INDEX idx_pianos_status_create_time ON pianos(status, create_time DESC);
CREATE INDEX idx_pianos_prefecture_city ON pianos(prefecture, city) WHERE status = 'active';
CREATE INDEX idx_pianos_kind_status ON pianos(kind, status);
CREATE INDEX idx_pianos_creator ON pianos(creator_user_id, create_time DESC);

CREATE TABLE IF NOT EXISTS piano_photos (
    id ulid PRIMARY KEY,
    piano_id ulid NOT NULL REFERENCES pianos(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    caption TEXT,
    display_order INT NOT NULL DEFAULT 0,
    uploader_user_id ulid REFERENCES users(id) ON DELETE SET NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_piano_photos_piano ON piano_photos(piano_id, display_order);

-- ============================================================================
-- Piano Posts (旧 performances + reviews 統合。1 訪問 = 1 投稿)
-- 投稿(post)は SNS の基本単位。rating は任意で付ける(評価を付けたい時だけ)。
-- メディア(動画/音声/画像)・コメントが付く。タイムラインに流れる対象。
-- ============================================================================

CREATE TYPE post_visibility AS ENUM ('public', 'private');
CREATE TYPE video_status AS ENUM ('queued', 'processing', 'ready', 'error');

CREATE TABLE IF NOT EXISTS piano_posts (
    id ulid PRIMARY KEY,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    piano_id ulid NOT NULL REFERENCES pianos(id) ON DELETE CASCADE,
    visit_time TIMESTAMPTZ NOT NULL,
    -- rating は必須(Google Maps 口コミ仕様)。最小投稿は rating だけでも可
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    body TEXT,
    -- 環境・楽器属性(任意、ピアノの「特徴メーター」表示用、評価には独立して集計される)
    ambient_noise SMALLINT CHECK (ambient_noise BETWEEN 1 AND 5),
    foot_traffic SMALLINT CHECK (foot_traffic BETWEEN 1 AND 5),
    resonance SMALLINT CHECK (resonance BETWEEN 1 AND 5),
    key_touch_weight SMALLINT CHECK (key_touch_weight BETWEEN 1 AND 5),
    tuning_quality SMALLINT CHECK (tuning_quality BETWEEN 1 AND 5),
    -- メディア
    video_uid TEXT UNIQUE,
    video_status video_status,
    video_duration_sec INT,
    video_thumbnail_url TEXT,
    audio_url TEXT,
    -- 公開範囲(private でも数値集計には反映する)
    visibility post_visibility NOT NULL DEFAULT 'public',
    comment_count INT NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_piano_posts_piano_create ON piano_posts(piano_id, create_time DESC);
CREATE INDEX idx_piano_posts_user_create ON piano_posts(user_id, create_time DESC);
CREATE INDEX idx_piano_posts_public_create ON piano_posts(create_time DESC) WHERE visibility = 'public';

CREATE TABLE IF NOT EXISTS piano_post_images (
    id ulid PRIMARY KEY,
    piano_post_id ulid NOT NULL REFERENCES piano_posts(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_piano_post_images_post ON piano_post_images(piano_post_id, display_order);

-- ============================================================================
-- Comments
-- ============================================================================

CREATE TABLE IF NOT EXISTS piano_post_comments (
    id ulid PRIMARY KEY,
    piano_post_id ulid NOT NULL REFERENCES piano_posts(id) ON DELETE CASCADE,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id ulid REFERENCES piano_post_comments(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_piano_post_comments_post ON piano_post_comments(piano_post_id, create_time);
CREATE INDEX idx_piano_post_comments_user ON piano_post_comments(user_id, create_time DESC);

CREATE TABLE IF NOT EXISTS piano_comments (
    id ulid PRIMARY KEY,
    piano_id ulid NOT NULL REFERENCES pianos(id) ON DELETE CASCADE,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id ulid REFERENCES piano_comments(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_piano_comments_piano ON piano_comments(piano_id, create_time);
CREATE INDEX idx_piano_comments_user ON piano_comments(user_id, create_time DESC);

-- ============================================================================
-- Piano Edits (公開編集ログ + Watch + 通知)
-- ============================================================================

CREATE TABLE IF NOT EXISTS piano_edits (
    id ulid PRIMARY KEY,
    piano_id ulid NOT NULL REFERENCES pianos(id) ON DELETE CASCADE,
    editor_user_id ulid REFERENCES users(id) ON DELETE SET NULL,
    operation piano_edit_operation NOT NULL,
    changes JSONB,
    summary TEXT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_piano_edits_piano_create ON piano_edits(piano_id, create_time DESC);
CREATE INDEX idx_piano_edits_editor_create ON piano_edits(editor_user_id, create_time DESC);

CREATE TABLE IF NOT EXISTS piano_watches (
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    piano_id ulid NOT NULL REFERENCES pianos(id) ON DELETE CASCADE,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, piano_id)
);

CREATE INDEX idx_piano_watches_piano ON piano_watches(piano_id);

-- 「行ってみたい」「行ったことある」「お気に入り」の3リストを単一テーブルで管理。
-- 同一 (user, piano) でも list_kind が違えば共存可能(複合主キー)。
-- 行ったことある(visited)は piano_post 作成時に usecase で UPSERT する想定。
CREATE TABLE IF NOT EXISTS piano_user_lists (
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    piano_id ulid NOT NULL REFERENCES pianos(id) ON DELETE CASCADE,
    list_kind piano_list_kind NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, piano_id, list_kind)
);

CREATE INDEX idx_piano_user_lists_user_kind ON piano_user_lists(user_id, list_kind, create_time DESC);
CREATE INDEX idx_piano_user_lists_piano_kind ON piano_user_lists(piano_id, list_kind);

-- ============================================================================
-- Notifications (MVPはDB蓄積 + クライアントポーリング、Phase1d で Redis pub/sub 追加)
-- ============================================================================

CREATE TYPE notification_type AS ENUM (
    'piano_edited',
    'piano_post_commented'
);

CREATE TABLE IF NOT EXISTS notifications (
    id ulid PRIMARY KEY,
    recipient_user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type notification_type NOT NULL,
    actor_user_id ulid REFERENCES users(id) ON DELETE CASCADE,
    resource_name TEXT NOT NULL DEFAULT '',
    read_time TIMESTAMPTZ,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_recipient_created ON notifications(recipient_user_id, create_time DESC);
CREATE INDEX idx_notifications_recipient_unread ON notifications(recipient_user_id) WHERE read_time IS NULL;

-- 同じ (recipient, type, actor, resource) の重複は防ぐ。NULL actor は空文字に正規化して比較。
CREATE UNIQUE INDEX idx_notifications_dedup
    ON notifications(recipient_user_id, type, COALESCE(actor_user_id::text, ''), resource_name);

-- ============================================================================
-- Reports
-- ============================================================================

CREATE TYPE report_target_type AS ENUM (
    'piano',
    'piano_post',
    'piano_post_comment',
    'piano_comment',
    'user'
);

CREATE TYPE report_status AS ENUM ('pending', 'reviewing', 'resolved', 'dismissed');

CREATE TYPE report_reason AS ENUM (
    'inappropriate',
    'spam',
    'copyright',
    'misinformation',
    'harassment',
    'privacy',
    'other'
);

CREATE TABLE IF NOT EXISTS reports (
    id ulid PRIMARY KEY,
    reporter_user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_type report_target_type NOT NULL,
    target_id ulid NOT NULL,
    reason report_reason NOT NULL,
    detail TEXT,
    status report_status NOT NULL DEFAULT 'pending',
    handler_user_id ulid REFERENCES users(id) ON DELETE SET NULL,
    handle_time TIMESTAMPTZ,
    resolution_note TEXT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(reporter_user_id, target_type, target_id)
);

CREATE INDEX idx_reports_target ON reports(target_type, target_id);
CREATE INDEX idx_reports_status_create ON reports(status, create_time);

-- ============================================================================
-- Content Hides (投稿・コメントの自動/手動非公開)
-- 通報閾値超過(MODERATION_AUTO_HIDE_THRESHOLD、MVPは5)で自動 hide、admin が手動 hide も可能。
-- 履歴・監査ログ兼用。「現在非公開か?」は EXISTS (... revoked_time IS NULL) で判定。
-- ============================================================================

CREATE TYPE hide_target_type AS ENUM (
    'piano_post',
    'piano_post_comment',
    'piano_comment'
);

CREATE TYPE hide_actor AS ENUM ('admin', 'auto_report_threshold');

CREATE TABLE IF NOT EXISTS content_hides (
    id ulid PRIMARY KEY,
    target_type hide_target_type NOT NULL,
    target_id ulid NOT NULL,
    actor hide_actor NOT NULL,
    reason TEXT,
    issued_by_user_id ulid REFERENCES users(id) ON DELETE SET NULL,  -- admin の場合のみセット
    revoked_time TIMESTAMPTZ,
    revoked_by_user_id ulid REFERENCES users(id) ON DELETE SET NULL,
    revoked_reason TEXT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_content_hides_target_active
  ON content_hides(target_type, target_id)
  WHERE revoked_time IS NULL;

CREATE INDEX idx_content_hides_target_create
  ON content_hides(target_type, target_id, create_time DESC);

-- ============================================================================
-- User Restrictions (悪質ユーザーの凍結・BAN)
-- 履歴と監査ログを兼ねる。1ユーザー複数回の制裁を時系列で残す。
-- 「現在制限中か?」は revoked_time IS NULL AND suspended_until > NOW() で判定。
-- ============================================================================

CREATE TABLE IF NOT EXISTS user_restrictions (
    id ulid PRIMARY KEY,
    user_id ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    suspended_until TIMESTAMPTZ NOT NULL,                   -- 'infinity' = 永久BAN、未来日 = 期限付き
    reason TEXT,
    issued_by_user_id ulid REFERENCES users(id) ON DELETE SET NULL,
    revoked_time TIMESTAMPTZ,                               -- NULL = 未解除
    revoked_by_user_id ulid REFERENCES users(id) ON DELETE SET NULL,
    revoked_reason TEXT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_restrictions_user_active
  ON user_restrictions(user_id, suspended_until DESC)
  WHERE revoked_time IS NULL;

CREATE INDEX idx_user_restrictions_user_create
  ON user_restrictions(user_id, create_time DESC);

-- ============================================================================
-- Triggers: denormalized aggregates
-- ============================================================================

-- piano_posts の変動を pianos の集計列に反映する。
-- rating は必須なので post_count == rating の母数。5属性は任意なので null-aware。
-- 集計は visibility に関わらず全件対象(private 投稿も数値だけはピアノに反映)。
CREATE FUNCTION piano_post_aggregates_sync() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE pianos SET
            post_count = post_count + 1,
            rating_sum = rating_sum + NEW.rating,
            ambient_noise_count = ambient_noise_count + (NEW.ambient_noise IS NOT NULL)::INT,
            ambient_noise_sum = ambient_noise_sum + COALESCE(NEW.ambient_noise, 0),
            foot_traffic_count = foot_traffic_count + (NEW.foot_traffic IS NOT NULL)::INT,
            foot_traffic_sum = foot_traffic_sum + COALESCE(NEW.foot_traffic, 0),
            resonance_count = resonance_count + (NEW.resonance IS NOT NULL)::INT,
            resonance_sum = resonance_sum + COALESCE(NEW.resonance, 0),
            key_touch_weight_count = key_touch_weight_count + (NEW.key_touch_weight IS NOT NULL)::INT,
            key_touch_weight_sum = key_touch_weight_sum + COALESCE(NEW.key_touch_weight, 0),
            tuning_quality_count = tuning_quality_count + (NEW.tuning_quality IS NOT NULL)::INT,
            tuning_quality_sum = tuning_quality_sum + COALESCE(NEW.tuning_quality, 0)
            WHERE id = NEW.piano_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE pianos SET
            post_count = GREATEST(post_count - 1, 0),
            rating_sum = GREATEST(rating_sum - OLD.rating, 0),
            ambient_noise_count = GREATEST(ambient_noise_count - (OLD.ambient_noise IS NOT NULL)::INT, 0),
            ambient_noise_sum = GREATEST(ambient_noise_sum - COALESCE(OLD.ambient_noise, 0), 0),
            foot_traffic_count = GREATEST(foot_traffic_count - (OLD.foot_traffic IS NOT NULL)::INT, 0),
            foot_traffic_sum = GREATEST(foot_traffic_sum - COALESCE(OLD.foot_traffic, 0), 0),
            resonance_count = GREATEST(resonance_count - (OLD.resonance IS NOT NULL)::INT, 0),
            resonance_sum = GREATEST(resonance_sum - COALESCE(OLD.resonance, 0), 0),
            key_touch_weight_count = GREATEST(key_touch_weight_count - (OLD.key_touch_weight IS NOT NULL)::INT, 0),
            key_touch_weight_sum = GREATEST(key_touch_weight_sum - COALESCE(OLD.key_touch_weight, 0), 0),
            tuning_quality_count = GREATEST(tuning_quality_count - (OLD.tuning_quality IS NOT NULL)::INT, 0),
            tuning_quality_sum = GREATEST(tuning_quality_sum - COALESCE(OLD.tuning_quality, 0), 0)
            WHERE id = OLD.piano_id;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE pianos SET
            rating_sum = rating_sum + (NEW.rating - OLD.rating),
            ambient_noise_count = ambient_noise_count
                + (NEW.ambient_noise IS NOT NULL)::INT
                - (OLD.ambient_noise IS NOT NULL)::INT,
            ambient_noise_sum = ambient_noise_sum
                + COALESCE(NEW.ambient_noise, 0)
                - COALESCE(OLD.ambient_noise, 0),
            foot_traffic_count = foot_traffic_count
                + (NEW.foot_traffic IS NOT NULL)::INT
                - (OLD.foot_traffic IS NOT NULL)::INT,
            foot_traffic_sum = foot_traffic_sum
                + COALESCE(NEW.foot_traffic, 0)
                - COALESCE(OLD.foot_traffic, 0),
            resonance_count = resonance_count
                + (NEW.resonance IS NOT NULL)::INT
                - (OLD.resonance IS NOT NULL)::INT,
            resonance_sum = resonance_sum
                + COALESCE(NEW.resonance, 0)
                - COALESCE(OLD.resonance, 0),
            key_touch_weight_count = key_touch_weight_count
                + (NEW.key_touch_weight IS NOT NULL)::INT
                - (OLD.key_touch_weight IS NOT NULL)::INT,
            key_touch_weight_sum = key_touch_weight_sum
                + COALESCE(NEW.key_touch_weight, 0)
                - COALESCE(OLD.key_touch_weight, 0),
            tuning_quality_count = tuning_quality_count
                + (NEW.tuning_quality IS NOT NULL)::INT
                - (OLD.tuning_quality IS NOT NULL)::INT,
            tuning_quality_sum = tuning_quality_sum
                + COALESCE(NEW.tuning_quality, 0)
                - COALESCE(OLD.tuning_quality, 0)
            WHERE id = NEW.piano_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_piano_posts_aggregates
AFTER INSERT OR UPDATE OR DELETE ON piano_posts
FOR EACH ROW EXECUTE FUNCTION piano_post_aggregates_sync();

CREATE FUNCTION piano_post_comment_counters_sync() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE piano_posts SET comment_count = comment_count + 1 WHERE id = NEW.piano_post_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE piano_posts SET comment_count = GREATEST(comment_count - 1, 0) WHERE id = OLD.piano_post_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_piano_post_comments_counters
AFTER INSERT OR DELETE ON piano_post_comments
FOR EACH ROW EXECUTE FUNCTION piano_post_comment_counters_sync();

-- users.post_count を piano_posts の変動に追従
CREATE FUNCTION user_post_counter_sync() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users SET post_count = post_count + 1 WHERE id = NEW.user_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users SET post_count = GREATEST(post_count - 1, 0) WHERE id = OLD.user_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_user_post_counter
AFTER INSERT OR DELETE ON piano_posts
FOR EACH ROW EXECUTE FUNCTION user_post_counter_sync();

-- users.edit_count を piano_edits の INSERT/DELETE に追従
CREATE FUNCTION user_edit_counter_sync() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.editor_user_id IS NOT NULL THEN
        UPDATE users SET edit_count = edit_count + 1 WHERE id = NEW.editor_user_id;
    ELSIF TG_OP = 'DELETE' AND OLD.editor_user_id IS NOT NULL THEN
        UPDATE users SET edit_count = GREATEST(edit_count - 1, 0) WHERE id = OLD.editor_user_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_user_edit_counter
AFTER INSERT OR DELETE ON piano_edits
FOR EACH ROW EXECUTE FUNCTION user_edit_counter_sync();

-- pianos の wishlist_count / visited_count / favorite_count を piano_user_lists の変動に追従
CREATE FUNCTION piano_user_list_counters_sync() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.list_kind = 'wishlist' THEN
            UPDATE pianos SET wishlist_count = wishlist_count + 1 WHERE id = NEW.piano_id;
        ELSIF NEW.list_kind = 'visited' THEN
            UPDATE pianos SET visited_count = visited_count + 1 WHERE id = NEW.piano_id;
        ELSIF NEW.list_kind = 'favorite' THEN
            UPDATE pianos SET favorite_count = favorite_count + 1 WHERE id = NEW.piano_id;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.list_kind = 'wishlist' THEN
            UPDATE pianos SET wishlist_count = GREATEST(wishlist_count - 1, 0) WHERE id = OLD.piano_id;
        ELSIF OLD.list_kind = 'visited' THEN
            UPDATE pianos SET visited_count = GREATEST(visited_count - 1, 0) WHERE id = OLD.piano_id;
        ELSIF OLD.list_kind = 'favorite' THEN
            UPDATE pianos SET favorite_count = GREATEST(favorite_count - 1, 0) WHERE id = OLD.piano_id;
        END IF;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_piano_user_list_counters
AFTER INSERT OR DELETE ON piano_user_lists
FOR EACH ROW EXECUTE FUNCTION piano_user_list_counters_sync();

-- users.report_received_count を reports.status='resolved' & target_type='user' の変動に追従
CREATE FUNCTION user_report_received_counter_sync() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.status = 'resolved' AND NEW.target_type = 'user' THEN
        UPDATE users SET report_received_count = report_received_count + 1 WHERE id = NEW.target_id;
    ELSIF TG_OP = 'UPDATE' AND NEW.target_type = 'user' THEN
        IF OLD.status <> 'resolved' AND NEW.status = 'resolved' THEN
            UPDATE users SET report_received_count = report_received_count + 1 WHERE id = NEW.target_id;
        ELSIF OLD.status = 'resolved' AND NEW.status <> 'resolved' THEN
            UPDATE users SET report_received_count = GREATEST(report_received_count - 1, 0) WHERE id = OLD.target_id;
        END IF;
    ELSIF TG_OP = 'DELETE' AND OLD.status = 'resolved' AND OLD.target_type = 'user' THEN
        UPDATE users SET report_received_count = GREATEST(report_received_count - 1, 0) WHERE id = OLD.target_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_user_report_received_counter
AFTER INSERT OR UPDATE OR DELETE ON reports
FOR EACH ROW EXECUTE FUNCTION user_report_received_counter_sync();

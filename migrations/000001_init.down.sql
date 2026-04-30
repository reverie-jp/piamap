-- drop triggers and functions
DROP TRIGGER IF EXISTS trg_user_report_received_counter ON reports;
DROP FUNCTION IF EXISTS user_report_received_counter_sync();

DROP TRIGGER IF EXISTS trg_piano_user_list_counters ON piano_user_lists;
DROP FUNCTION IF EXISTS piano_user_list_counters_sync();

DROP TRIGGER IF EXISTS trg_user_edit_counter ON piano_edits;
DROP FUNCTION IF EXISTS user_edit_counter_sync();

DROP TRIGGER IF EXISTS trg_user_post_counter ON piano_posts;
DROP FUNCTION IF EXISTS user_post_counter_sync();

DROP TRIGGER IF EXISTS trg_piano_post_comments_counters ON piano_post_comments;
DROP FUNCTION IF EXISTS piano_post_comment_counters_sync();

DROP TRIGGER IF EXISTS trg_piano_posts_aggregates ON piano_posts;
DROP FUNCTION IF EXISTS piano_post_aggregates_sync();

-- drop tables
DROP TABLE IF EXISTS user_restrictions;
DROP TABLE IF EXISTS content_hides;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS piano_user_lists;
DROP TABLE IF EXISTS piano_watches;
DROP TABLE IF EXISTS piano_edits;
DROP TABLE IF EXISTS piano_comments;
DROP TABLE IF EXISTS piano_post_comments;
DROP TABLE IF EXISTS piano_post_images;
DROP TABLE IF EXISTS piano_posts;
DROP TABLE IF EXISTS piano_photos;
DROP TABLE IF EXISTS pianos;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_auth_providers;
DROP TABLE IF EXISTS users;

-- drop types
DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS hide_actor;
DROP TYPE IF EXISTS hide_target_type;
DROP TYPE IF EXISTS report_reason;
DROP TYPE IF EXISTS report_status;
DROP TYPE IF EXISTS report_target_type;
DROP TYPE IF EXISTS video_status;
DROP TYPE IF EXISTS post_visibility;
DROP TYPE IF EXISTS piano_edit_operation;
DROP TYPE IF EXISTS piano_availability;
DROP TYPE IF EXISTS piano_list_kind;
DROP TYPE IF EXISTS piano_kind;
DROP TYPE IF EXISTS piano_type;
DROP TYPE IF EXISTS piano_status;
DROP TYPE IF EXISTS auth_provider;

-- drop domains
DROP DOMAIN IF EXISTS ulid;

-- PostGIS extension は他データベースで共有される可能性があるため残置

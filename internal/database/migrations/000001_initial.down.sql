-- Rollback initial schema
-- This will drop all tables created in the up migration

PRAGMA foreign_keys=OFF;

DROP INDEX IF EXISTS idx_brags_job_title_id;
DROP INDEX IF EXISTS idx_brags_category_id;
DROP INDEX IF EXISTS idx_brags_owner_id;
DROP INDEX IF EXISTS idx_tags_name;
DROP INDEX IF EXISTS idx_tags_owner_id;
DROP INDEX IF EXISTS idx_brag_tags_tag_id;
DROP INDEX IF EXISTS idx_brag_tags_brag_id;
DROP INDEX IF EXISTS idx_job_titles_user_active;

DROP TABLE IF EXISTS brag_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS brags;
DROP TABLE IF EXISTS job_titles;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

PRAGMA foreign_keys=ON;

-- Migration 004: Add unique constraint for active job titles
-- Ensures only one active job title per user (where end_date IS NULL)

CREATE UNIQUE INDEX IF NOT EXISTS idx_job_titles_user_active 
ON job_titles(user_id) 
WHERE end_date IS NULL;

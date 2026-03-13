-- name: CountBragsByJobTitle :one
SELECT COUNT(*) FROM brags WHERE job_title_id = ?;

-- name: GetJobTitle :one
SELECT * FROM job_titles WHERE id = ? LIMIT 1;

-- name: GetActiveJobTitle :one
SELECT * FROM job_titles 
WHERE user_id = ? AND end_date IS NULL 
ORDER BY start_date DESC, created_at DESC 
LIMIT 1;

-- name: GetJobTitleByName :one
SELECT * FROM job_titles 
WHERE user_id = ? AND title = ? 
ORDER BY start_date DESC, created_at DESC 
LIMIT 1;

-- name: ListJobTitlesByUser :many
SELECT * FROM job_titles
WHERE user_id = ?
ORDER BY CASE WHEN start_date IS NULL THEN 1 ELSE 0 END, start_date DESC, created_at DESC;

-- name: CreateJobTitle :one
INSERT INTO job_titles (user_id, title, company, start_date, end_date, created_at)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateJobTitle :one
UPDATE job_titles
SET title = ?, company = ?, start_date = ?, end_date = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteJobTitle :exec
DELETE FROM job_titles WHERE id = ?;

-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (name, email, job_title, company, language, created_at)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET name = ?, email = ?, job_title = ?, company = ?, language = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

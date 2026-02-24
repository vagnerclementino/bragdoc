-- name: CountBragsByPosition :one
SELECT COUNT(*) FROM brags WHERE position_id = ?;

-- name: GetPosition :one
SELECT * FROM positions WHERE id = ? LIMIT 1;

-- name: ListPositionsByUser :many
SELECT * FROM positions
WHERE user_id = ?
ORDER BY CASE WHEN start_date IS NULL THEN 1 ELSE 0 END, start_date DESC, created_at DESC;

-- name: CreatePosition :one
INSERT INTO positions (user_id, title, company, start_date, end_date, created_at)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdatePosition :one
UPDATE positions
SET title = ?, company = ?, start_date = ?, end_date = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeletePosition :exec
DELETE FROM positions WHERE id = ?;

-- name: GetBrag :one
SELECT * FROM brags WHERE id = ? LIMIT 1;

-- name: ListBragsByUser :many
SELECT * FROM brags WHERE owner_id = ? ORDER BY created_at;

-- name: ListBragsByCategory :many
SELECT * FROM brags WHERE owner_id = ? AND category_id = ? ORDER BY created_at;

-- name: CreateBrag :one
INSERT INTO brags (owner_id, title, description, category_id, position_id, created_at)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateBrag :one
UPDATE brags 
SET title = ?, description = ?, category_id = ?, position_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteBrag :exec
DELETE FROM brags WHERE id = ?;

-- name: SearchBragsByTags :many
SELECT DISTINCT b.* FROM brags b
JOIN brag_tags bt ON b.id = bt.brag_id
JOIN tags t ON bt.tag_id = t.id
WHERE b.owner_id = ? AND t.name IN (sqlc.slice('tag_names'))
ORDER BY b.created_at DESC;

-- name: GetTag :one
SELECT * FROM tags WHERE id = ? LIMIT 1;

-- name: GetTagByName :one
SELECT * FROM tags WHERE owner_id = ? AND name = ? LIMIT 1;

-- name: ListTagsByUser :many
SELECT * FROM tags WHERE owner_id = ? ORDER BY name ASC;

-- name: ListTagsByBrag :many
SELECT t.* FROM tags t
JOIN brag_tags bt ON t.id = bt.tag_id
WHERE bt.brag_id = ?
ORDER BY t.name ASC;

-- name: CreateTag :one
INSERT INTO tags (name, owner_id, created_at)
VALUES (?, ?, CURRENT_TIMESTAMP)
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags WHERE id = ?;

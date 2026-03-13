-- name: CountBragsByCategory :one
SELECT COUNT(*) FROM brags WHERE category_id = ?;

-- name: CreateCategory :one
INSERT INTO categories (name, description, created_at, updated_at)
VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = ?;

-- name: GetCategory :one
SELECT * FROM categories WHERE id = ? LIMIT 1;

-- name: GetCategoryByName :one
SELECT * FROM categories WHERE name = ? LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY name;

-- name: UpdateCategory :one
UPDATE categories 
SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;
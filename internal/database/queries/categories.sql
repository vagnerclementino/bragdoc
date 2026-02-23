-- name: GetCategory :one
SELECT * FROM categories WHERE id = ? LIMIT 1;

-- name: GetCategoryByName :one
SELECT * FROM categories WHERE name = ? LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY name;

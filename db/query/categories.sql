-- name: CreateCategory :one
INSERT INTO categories (name, type)
VALUES ($1, $2)
RETURNING id, name, type, created_at;

-- name: GetCategory :one
SELECT id, name, type, created_at FROM categories
WHERE id = $1;

-- name: ListCategories :many
SELECT id, name, type, created_at FROM categories
ORDER BY id;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, type = $3
WHERE id = $1
RETURNING id, name, type, created_at;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;
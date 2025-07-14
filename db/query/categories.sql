-- name: CreateCategory :one
INSERT INTO categories (name, type, account_id, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories WHERE id = $1;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, type = $3, account_id = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: GetCategoryIDByName :one
SELECT id FROM categories WHERE name = $1;
-- name: CreateProduct :one
INSERT INTO products (name, price, category_id, value)
VALUES ($1, $2, $3, $4)
RETURNING id, name, price, category_id, value, created_at;

-- name: GetProduct :one
SELECT
  p.id,
  p.name,
  p.price,
  c.name AS category_name,
  c.type AS category_type,
  p.value
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
WHERE
  p.id = $1;

-- name: ListProducts :many
SELECT id, name, price, category_id, value, created_at FROM products
ORDER BY id;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, price = $3, category_id = $4, value = $5
WHERE id = $1
RETURNING id, name, price, category_id, value, created_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: ListProductsByCategoryID :many
SELECT
  p.id,
  p.name,
  p.price,
  p.category_id,
  p.value,
  p.created_at
FROM
  products AS p
WHERE
  p.category_id = $1
ORDER BY p.id;
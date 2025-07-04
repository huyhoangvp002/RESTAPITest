
-- name: CreateProduct :one
INSERT INTO products (name, price, category_id)
VALUES ($1, $2, $3)
RETURNING id, name, price, category_id, created_at;

-- name: GetProduct :one
SELECT id, name, price, category_id, created_at FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, price, category_id, created_at FROM products
ORDER BY id;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, price = $3, category_id = $4
WHERE id = $1
RETURNING id, name, price, category_id, created_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: GetProductWithCategory :one
SELECT
  p.id,
  p.name,
  p.price,
  p.category_id,
  c.name AS category_name,
  p.created_at
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
WHERE
  p.id = $1;
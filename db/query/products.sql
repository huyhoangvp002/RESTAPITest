-- name: CreateProduct :one
INSERT INTO products (
  account_id, category_id, name, price, discount_price, stock_quantity
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetProduct :one
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  c.name AS category_name,
  c.type AS category_type,
  p.stock_quantity,
  p.account_id,
  p.created_at
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
WHERE
  p.id = $1;

-- name: ListProducts :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  c.name AS category_name,
  c.type AS category_type,
  p.stock_quantity
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
ORDER BY p.id
LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
UPDATE products
SET
  price = $2,
  stock_quantity = $3
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: ListProductsByCategoryID :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.stock_quantity,
  p.account_id,
  p.created_at
FROM
  products AS p
WHERE
  p.category_id = $1
ORDER BY p.id
LIMIT $2 OFFSET $3;

-- name: ListProductsByMaxPrice :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.stock_quantity,
  p.account_id,
  p.created_at,
  c.name AS category_name
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
WHERE
  p.discount_price < $1
ORDER BY p.discount_price ASC;

-- name: UpdateDiscountPrice :exec
UPDATE products
SET discount_price = $2 WHERE id = $1;

-- name: GetPriceByID :one
SELECT price FROM products WHERE id = $1;

-- name: ListProductByAccountID :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.stock_quantity,
  p.account_id,
  p.created_at
FROM
  products AS p
WHERE
  p.account_id = $1
ORDER BY p.id
LIMIT $2 OFFSET $3;

-- name: GetProdIDByAccountID :one
SELECT p.id FROM products AS p WHERE p.account_id = $1;

-- name: GetAccountIDbyProductID :one
SELECT account_id FROM products WHERE id = $1;

-- name: SearchProductsByName :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.stock_quantity,
  c.name AS category_name
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
WHERE
  p.name ILIKE '%' || $1 || '%'
ORDER BY p.id
LIMIT $2 OFFSET $3;

-- name: UpdateProductStockByID :exec
UPDATE products
SET stock_quantity = $1
WHERE id = $2;

-- name: GetDiscountPriceByID :one
SELECT discount_price FROM products WHERE id = $1;

-- name: GetStockByID :one
SELECT stock_quantity FROM products WHERE id = $1;
-- name: CreateProduct :one
INSERT INTO products (
    name, price, discount_price, category_id, value, account_id, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, name, price, discount_price, category_id, value, account_id, created_at;

-- name: GetProduct :one
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  c.name AS category_name,
  c.type AS category_type,
  p.value,
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
  p.value,
  p.account_id,
  p.created_at
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
  value = $3
WHERE id = $1
RETURNING id, name, price, discount_price, value, account_id, category_id, created_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: ListProductsByCategoryID :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.value,
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
  p.value,
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
SET
  discount_price = $2
WHERE
  id = $1;

-- name: GetPriceByID :one
SELECT price
FROM products
WHERE id = $1;

-- name: ListProductByAccountID :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.value,
  p.account_id,
  p.created_at
FROM
  products AS p
WHERE
  p.account_id = $1
ORDER BY p.id
LIMIT $2 OFFSET $3;

-- name: GetProdIDByAccountID :one
SELECT
  p.id
FROM
  products AS p
WHERE
  p.account_id = $1;

-- name: GetAccountIDbyProductID :one
SELECT account_id FROM products WHERE id = $1;

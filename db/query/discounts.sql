-- name: CreateDiscount :one
INSERT INTO discounts (
  discount_value,
  product_id,
  customer_id
) VALUES (
  $1, $2, $3
)
RETURNING id, discount_value, product_id, customer_id, created_at;

-- name: GetDiscount :one
SELECT
  id,
  discount_value,
  product_id,
  customer_id,
  created_at
FROM
  discounts
WHERE
  id = $1;

-- name: ListDiscounts :many
SELECT
  id,
  discount_value,
  product_id,
  customer_id,
  created_at
FROM
  discounts
ORDER BY
  id;

-- name: UpdateDiscount :one
UPDATE discounts
SET
  discount_value = $2
WHERE
  id = $1
RETURNING id, discount_value, product_id, customer_id, created_at;

-- name: DeleteDiscount :exec
DELETE FROM discounts
WHERE id = $1;

-- name: ListDiscountsByCustomerID :many
SELECT
  id,
  discount_value,
  product_id,
  customer_id,
  created_at
FROM
  discounts
WHERE
  customer_id = $1
ORDER BY
  id;

-- name: GetProductIDByCustomerID :one
SELECT
  product_id
FROM
  discounts
WHERE
  customer_id = $1;

-- name: CreateDiscount :one
INSERT INTO discounts (
  discount_value,
  account_id,
  product_id,
  created_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, discount_value, account_id, product_id, created_at;

-- name: GetDiscount :one
SELECT
  id,
  discount_value,
  account_id,
  product_id,
  created_at
FROM
  discounts
WHERE
  id = $1;

-- name: ListDiscounts :many
SELECT
  id,
  discount_value,
  account_id,
  product_id,
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
RETURNING id, discount_value, account_id, product_id, created_at;

-- name: DeleteDiscount :exec
DELETE FROM discounts
WHERE id = $1;

-- name: ListDiscountsByAccountID :many
SELECT
  id,
  discount_value,
  account_id,
  product_id,
  created_at
FROM
  discounts
WHERE
  account_id = $1
ORDER BY
  id;

-- name: GetProductIDByAccountID :one
SELECT
  product_id
FROM
  discounts
WHERE
  account_id = $1;

--
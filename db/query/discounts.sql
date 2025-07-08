-- name: CreateDiscount :one
INSERT INTO discounts (
  discount_value,
  product_id
) VALUES (
  $1, $2
)
RETURNING id, discount_value, product_id;

-- name: GetDiscount :one
SELECT
  id,
  discount_value,
  product_id
FROM
  discounts
WHERE
  id = $1;

-- name: ListDiscounts :many
SELECT
  id,
  discount_value,
  product_id
FROM
  discounts
ORDER BY
  id;

-- name: UpdateDiscount :one
UPDATE discounts
SET
  discount_value = $2,
  product_id = $3
WHERE
  id = $1
RETURNING id, discount_value, product_id;

-- name: DeleteDiscount :exec
DELETE FROM discounts
WHERE id = $1;

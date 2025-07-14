-- name: CreateDiscount :one
INSERT INTO discounts (
  discount_value,
  account_id,
  product_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetDiscount :one
SELECT * FROM discounts WHERE id = $1;

-- name: ListDiscounts :many
SELECT * FROM discounts ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateDiscount :one
UPDATE discounts
SET discount_value = $2
WHERE id = $1
RETURNING *;

-- name: DeleteDiscount :exec
DELETE FROM discounts WHERE id = $1;

-- name: ListDiscountsByAccountID :many
SELECT * FROM discounts WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: GetProductIDByAccountID :one
SELECT product_id FROM discounts WHERE account_id = $1;

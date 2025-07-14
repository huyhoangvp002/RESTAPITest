-- name: CreateCartItem :one
INSERT INTO cart_items (
  quantity,
  account_id,
  product_id,
  created_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE id = $1;

-- name: GetAccountIDByCartItemID :one
SELECT account_id FROM cart_items WHERE id = $1;

-- name: ListCartItemsByAccountID :many
SELECT
  p.name AS product_name,
  c.id,
  c.quantity
FROM
  cart_items AS c
JOIN
  products AS p ON c.product_id = p.id
WHERE
  c.account_id = $1
ORDER BY
  c.id
LIMIT $2 OFFSET $3;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM cart_items WHERE id = $1;

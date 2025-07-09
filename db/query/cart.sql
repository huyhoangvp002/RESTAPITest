-- name: CreateCart :one
INSERT INTO cart (
  value,
  account_id,
  product_id,
  created_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetCart :one
SELECT * FROM cart
WHERE id = $1;

-- name: ListCarts :many
SELECT * FROM cart
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateCartValue :one
UPDATE cart
SET
  value = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCart :exec
DELETE FROM cart
WHERE id = $1;

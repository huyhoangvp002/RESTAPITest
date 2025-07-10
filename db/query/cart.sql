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

-- name: GetAccountIDByCartID :one
SELECT account_id
FROM cart
WHERE id = $1;

-- name: ListCartByAccountID :many
SELECT
 p.name AS product_name,
  c.id,
  c.value
 
FROM
  cart AS c
JOIN
  products AS p ON c.product_id = p.id
WHERE
  c.account_id = $1
ORDER BY
  c.id
LIMIT $2 OFFSET $3;

-- name: UpdateCartValue :one
UPDATE cart
SET
  value = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCart :exec
DELETE FROM cart
WHERE id = $1;



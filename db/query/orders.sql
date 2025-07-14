-- name: CreateOrder :one
INSERT INTO orders (buyer_id, seller_id, total_price, cod, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateOrderStatus :one
UPDATE orders SET status = $2 WHERE id = $1 RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;
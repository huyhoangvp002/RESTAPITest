-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_each)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetOrderItem :one
SELECT * FROM order_items WHERE id = $1;

-- name: ListOrderItemsByOrderID :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: UpdateOrderItemQuantity :one
UPDATE order_items SET quantity = $2 WHERE id = $1 RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items WHERE id = $1;
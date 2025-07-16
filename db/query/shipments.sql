-- name: CreateShipment :one
INSERT INTO shipments (order_id, shipment_code, fee, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetShipment :one
SELECT * FROM shipments WHERE id = $1;

-- name: ListShipments :many
SELECT * FROM shipments ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateShipmentStatus :one
UPDATE shipments SET status = $2, updated_at = $3 WHERE id = $1 RETURNING *;

-- name: DeleteShipment :exec
DELETE FROM shipments WHERE id = $1;
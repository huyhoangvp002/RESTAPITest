-- name: CreateCustomer :one
INSERT INTO customers (
    name,
    account_id,
    email
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE id = $1;

-- name: ListCustomers :many
SELECT * FROM customers
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateCustomer :one
UPDATE customers
SET
    name = $2,
    email = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;

-- name: GetCustomerIDByUsername :one
SELECT c.id
FROM customers AS c
JOIN accounts AS a ON c.account_id = a.id
WHERE a.username = $1;
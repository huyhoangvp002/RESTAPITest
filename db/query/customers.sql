-- name: CreateCustomer :one
INSERT INTO customers (
    name,
    accounts_id,
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
SELECT
  cu.id
FROM
  customers AS cu
JOIN
  accounts AS ac ON cu.accounts_id = ac.id
WHERE
  ac.username = $1;
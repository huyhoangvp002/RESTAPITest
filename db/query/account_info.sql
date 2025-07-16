-- name: CreateAccountInfo :one
INSERT INTO account_info (
  name,
  email,
  phone_number,
  address,
  account_id,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetAccountInfo :one
SELECT * FROM account_info WHERE id = $1;

-- name: GetAccountID :one
SELECT account_id FROM account_info WHERE id = $1;

-- name: ListAccountInfos :many
SELECT * FROM account_info ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateAccountInfo :one
UPDATE account_info
SET
  name = $2,
  email = $3,
  phone_number = $4,
  address = $5,
  account_id = $6,
  updated_at = $7
WHERE id = $1
RETURNING *;

-- name: UpdateAccountInfoName :exec
UPDATE account_info SET name = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateAccountInfoEmail :exec
UPDATE account_info SET email = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateAccountInfoPhoneNumber :exec
UPDATE account_info SET phone_number = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateAccountInfoAddress :exec
UPDATE account_info SET address = $2, updated_at = NOW() WHERE id = $1;

-- name: GetNameForShipment :one
SELECT name FROM account_info WHERE account_id = $1;

-- name: GetPhoneForShipment :one
SELECT phone_number FROM account_info WHERE account_id = $1;

-- name: GetAddressForShipment :one
SELECT address FROM account_info WHERE account_id = $1;


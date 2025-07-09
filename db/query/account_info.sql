-- name: CreateAccountInfo :one
INSERT INTO account_info (
  name,
  email,
  phone_number,
  address,
  account_id,
  created_at,
  update_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetAccountInfo :one
SELECT * FROM account_info
WHERE id = $1;

-- name: ListAccountInfos :many
SELECT * FROM account_info
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateAccountInfo :one
UPDATE account_info
SET
  name = $2,
  email = $3,
  phone_number = $4,
  address = $5,
  account_id = $6,
  update_at = $7
WHERE id = $1
RETURNING *;

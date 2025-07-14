-- name: CreateAccount :one
INSERT INTO accounts (username, hash_password, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetAccountByUsername :one
SELECT * FROM accounts WHERE username = $1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET username = $2,
    hash_password = $3,
    role = $4
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;

-- name: GetIDByUserName :one
SELECT id FROM accounts WHERE username = $1;

-- name: UpdateRole :one
UPDATE accounts SET role = $2 WHERE id = $1 RETURNING *;

-- name: GetAccountIDByUsername :one
SELECT id FROM accounts WHERE username = $1;
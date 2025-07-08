-- name: CreateAccount :one
INSERT INTO accounts (username, hash_password, role)
VALUES ($1, $2, $3)
RETURNING id, username, role;

-- name: GetAccountByID :one
SELECT id, username, role
FROM accounts
WHERE id = $1;

-- name: GetAccountByUsername :one
SELECT username, hash_password, role
FROM accounts
WHERE username = $1;

-- name: ListAccounts :many
SELECT id, username, role
FROM accounts
ORDER BY id 
LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET username = $2,
    hash_password = $3,
    role = $4
WHERE id = $1
RETURNING id, username, role;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;


-- name: GetIDByUserName :one
SELECT id FROM accounts WHERE username =$1;

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: accounts.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (username, hash_password, role)
VALUES ($1, $2, $3)
RETURNING id, username, role
`

type CreateAccountParams struct {
	Username     string `json:"username"`
	HashPassword string `json:"hash_password"`
	Role         string `json:"role"`
}

type CreateAccountRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (CreateAccountRow, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Username, arg.HashPassword, arg.Role)
	var i CreateAccountRow
	err := row.Scan(&i.ID, &i.Username, &i.Role)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccountByID = `-- name: GetAccountByID :one
SELECT id, username, role
FROM accounts
WHERE id = $1
`

type GetAccountByIDRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (q *Queries) GetAccountByID(ctx context.Context, id int64) (GetAccountByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getAccountByID, id)
	var i GetAccountByIDRow
	err := row.Scan(&i.ID, &i.Username, &i.Role)
	return i, err
}

const getAccountByUsername = `-- name: GetAccountByUsername :one
SELECT username, hash_password, role
FROM accounts
WHERE username = $1
`

type GetAccountByUsernameRow struct {
	Username     string `json:"username"`
	HashPassword string `json:"hash_password"`
	Role         string `json:"role"`
}

func (q *Queries) GetAccountByUsername(ctx context.Context, username string) (GetAccountByUsernameRow, error) {
	row := q.db.QueryRowContext(ctx, getAccountByUsername, username)
	var i GetAccountByUsernameRow
	err := row.Scan(&i.Username, &i.HashPassword, &i.Role)
	return i, err
}

const getAccountIDByUsername = `-- name: GetAccountIDByUsername :one
SELECT id
FROM accounts
WHERE username = $1
`

func (q *Queries) GetAccountIDByUsername(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getAccountIDByUsername, username)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getIDByUserName = `-- name: GetIDByUserName :one
SELECT id FROM accounts WHERE username = $1
`

func (q *Queries) GetIDByUserName(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getIDByUserName, username)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, username, role
FROM accounts
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListAccountsRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]ListAccountsRow, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAccountsRow{}
	for rows.Next() {
		var i ListAccountsRow
		if err := rows.Scan(&i.ID, &i.Username, &i.Role); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET username = $2,
    hash_password = $3,
    role = $4
WHERE id = $1
RETURNING id, username, role
`

type UpdateAccountParams struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	HashPassword string `json:"hash_password"`
	Role         string `json:"role"`
}

type UpdateAccountRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (UpdateAccountRow, error) {
	row := q.db.QueryRowContext(ctx, updateAccount,
		arg.ID,
		arg.Username,
		arg.HashPassword,
		arg.Role,
	)
	var i UpdateAccountRow
	err := row.Scan(&i.ID, &i.Username, &i.Role)
	return i, err
}

const updateRole = `-- name: UpdateRole :one
UPDATE accounts
SET role = $2
WHERE id = $1
RETURNING id, username, role
`

type UpdateRoleParams struct {
	ID   int64  `json:"id"`
	Role string `json:"role"`
}

type UpdateRoleRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) (UpdateRoleRow, error) {
	row := q.db.QueryRowContext(ctx, updateRole, arg.ID, arg.Role)
	var i UpdateRoleRow
	err := row.Scan(&i.ID, &i.Username, &i.Role)
	return i, err
}

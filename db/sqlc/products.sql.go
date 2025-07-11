// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: products.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
    name, price, discount_price, category_id, value, account_id, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, name, price, discount_price, category_id, value, account_id, created_at
`

type CreateProductParams struct {
	Name          string        `json:"name"`
	Price         int32         `json:"price"`
	DiscountPrice int32         `json:"discount_price"`
	CategoryID    sql.NullInt32 `json:"category_id"`
	Value         int32         `json:"value"`
	AccountID     sql.NullInt32 `json:"account_id"`
	CreatedAt     time.Time     `json:"created_at"`
}

type CreateProductRow struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Price         int32         `json:"price"`
	DiscountPrice int32         `json:"discount_price"`
	CategoryID    sql.NullInt32 `json:"category_id"`
	Value         int32         `json:"value"`
	AccountID     sql.NullInt32 `json:"account_id"`
	CreatedAt     time.Time     `json:"created_at"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (CreateProductRow, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.Name,
		arg.Price,
		arg.DiscountPrice,
		arg.CategoryID,
		arg.Value,
		arg.AccountID,
		arg.CreatedAt,
	)
	var i CreateProductRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Price,
		&i.DiscountPrice,
		&i.CategoryID,
		&i.Value,
		&i.AccountID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getAccountIDbyProductID = `-- name: GetAccountIDbyProductID :one
SELECT account_id FROM products WHERE id = $1
`

func (q *Queries) GetAccountIDbyProductID(ctx context.Context, id int64) (sql.NullInt32, error) {
	row := q.db.QueryRowContext(ctx, getAccountIDbyProductID, id)
	var account_id sql.NullInt32
	err := row.Scan(&account_id)
	return account_id, err
}

const getPriceByID = `-- name: GetPriceByID :one
SELECT price
FROM products
WHERE id = $1
`

func (q *Queries) GetPriceByID(ctx context.Context, id int64) (int32, error) {
	row := q.db.QueryRowContext(ctx, getPriceByID, id)
	var price int32
	err := row.Scan(&price)
	return price, err
}

const getProdIDByAccountID = `-- name: GetProdIDByAccountID :one
SELECT
  p.id
FROM
  products AS p
WHERE
  p.account_id = $1
`

func (q *Queries) GetProdIDByAccountID(ctx context.Context, accountID sql.NullInt32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getProdIDByAccountID, accountID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getProduct = `-- name: GetProduct :one
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  c.name AS category_name,
  c.type AS category_type,
  p.value,
  p.account_id,
  p.created_at
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
WHERE
  p.id = $1
`

type GetProductRow struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Price         int32         `json:"price"`
	DiscountPrice int32         `json:"discount_price"`
	CategoryName  string        `json:"category_name"`
	CategoryType  string        `json:"category_type"`
	Value         int32         `json:"value"`
	AccountID     sql.NullInt32 `json:"account_id"`
	CreatedAt     time.Time     `json:"created_at"`
}

func (q *Queries) GetProduct(ctx context.Context, id int64) (GetProductRow, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i GetProductRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Price,
		&i.DiscountPrice,
		&i.CategoryName,
		&i.CategoryType,
		&i.Value,
		&i.AccountID,
		&i.CreatedAt,
	)
	return i, err
}

const listProductByAccountID = `-- name: ListProductByAccountID :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.value,
  p.account_id,
  p.created_at
FROM
  products AS p
WHERE
  p.account_id = $1
ORDER BY p.id
LIMIT $2 OFFSET $3
`

type ListProductByAccountIDParams struct {
	AccountID sql.NullInt32 `json:"account_id"`
	Limit     int32         `json:"limit"`
	Offset    int32         `json:"offset"`
}

type ListProductByAccountIDRow struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Price         int32         `json:"price"`
	DiscountPrice int32         `json:"discount_price"`
	Value         int32         `json:"value"`
	AccountID     sql.NullInt32 `json:"account_id"`
	CreatedAt     time.Time     `json:"created_at"`
}

func (q *Queries) ListProductByAccountID(ctx context.Context, arg ListProductByAccountIDParams) ([]ListProductByAccountIDRow, error) {
	rows, err := q.db.QueryContext(ctx, listProductByAccountID, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListProductByAccountIDRow{}
	for rows.Next() {
		var i ListProductByAccountIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.Value,
			&i.AccountID,
			&i.CreatedAt,
		); err != nil {
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

const listProducts = `-- name: ListProducts :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  c.name AS category_name,
  c.type AS category_type,
  p.value,
  p.account_id,
  p.created_at
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
ORDER BY p.id
LIMIT $1 OFFSET $2
`

type ListProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListProductsRow struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Price         int32         `json:"price"`
	DiscountPrice int32         `json:"discount_price"`
	CategoryName  string        `json:"category_name"`
	CategoryType  string        `json:"category_type"`
	Value         int32         `json:"value"`
	AccountID     sql.NullInt32 `json:"account_id"`
	CreatedAt     time.Time     `json:"created_at"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]ListProductsRow, error) {
	rows, err := q.db.QueryContext(ctx, listProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListProductsRow{}
	for rows.Next() {
		var i ListProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.CategoryName,
			&i.CategoryType,
			&i.Value,
			&i.AccountID,
			&i.CreatedAt,
		); err != nil {
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

const listProductsByCategoryID = `-- name: ListProductsByCategoryID :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.value,
  p.account_id,
  p.created_at
FROM
  products AS p
WHERE
  p.category_id = $1
ORDER BY p.id
LIMIT $2 OFFSET $3
`

type ListProductsByCategoryIDParams struct {
	CategoryID sql.NullInt32 `json:"category_id"`
	Limit      int32         `json:"limit"`
	Offset     int32         `json:"offset"`
}

type ListProductsByCategoryIDRow struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Price         int32         `json:"price"`
	DiscountPrice int32         `json:"discount_price"`
	Value         int32         `json:"value"`
	AccountID     sql.NullInt32 `json:"account_id"`
	CreatedAt     time.Time     `json:"created_at"`
}

func (q *Queries) ListProductsByCategoryID(ctx context.Context, arg ListProductsByCategoryIDParams) ([]ListProductsByCategoryIDRow, error) {
	rows, err := q.db.QueryContext(ctx, listProductsByCategoryID, arg.CategoryID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListProductsByCategoryIDRow{}
	for rows.Next() {
		var i ListProductsByCategoryIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.Value,
			&i.AccountID,
			&i.CreatedAt,
		); err != nil {
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

const listProductsByMaxPrice = `-- name: ListProductsByMaxPrice :many
SELECT
  p.id,
  p.name,
  p.price,
  p.discount_price,
  p.value,
  p.account_id,
  p.created_at,
  c.name AS category_name
FROM
  products AS p
JOIN
  categories AS c ON p.category_id = c.id
WHERE
  p.discount_price < $1
ORDER BY p.discount_price ASC
`

type ListProductsByMaxPriceRow struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Price         int32         `json:"price"`
	DiscountPrice int32         `json:"discount_price"`
	Value         int32         `json:"value"`
	AccountID     sql.NullInt32 `json:"account_id"`
	CreatedAt     time.Time     `json:"created_at"`
	CategoryName  string        `json:"category_name"`
}

func (q *Queries) ListProductsByMaxPrice(ctx context.Context, discountPrice int32) ([]ListProductsByMaxPriceRow, error) {
	rows, err := q.db.QueryContext(ctx, listProductsByMaxPrice, discountPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListProductsByMaxPriceRow{}
	for rows.Next() {
		var i ListProductsByMaxPriceRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.Value,
			&i.AccountID,
			&i.CreatedAt,
			&i.CategoryName,
		); err != nil {
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

const updateDiscountPrice = `-- name: UpdateDiscountPrice :exec
UPDATE products
SET
  discount_price = $2
WHERE
  id = $1
`

type UpdateDiscountPriceParams struct {
	ID            int64 `json:"id"`
	DiscountPrice int32 `json:"discount_price"`
}

func (q *Queries) UpdateDiscountPrice(ctx context.Context, arg UpdateDiscountPriceParams) error {
	_, err := q.db.ExecContext(ctx, updateDiscountPrice, arg.ID, arg.DiscountPrice)
	return err
}

const updateProduct = `-- name: UpdateProduct :one
UPDATE products
SET
  price = $2,
  value = $3
WHERE id = $1
RETURNING id, name, price, discount_price, value, account_id, category_id, created_at
`

type UpdateProductParams struct {
	ID    int64 `json:"id"`
	Price int32 `json:"price"`
	Value int32 `json:"value"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct, arg.ID, arg.Price, arg.Value)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Price,
		&i.DiscountPrice,
		&i.Value,
		&i.AccountID,
		&i.CategoryID,
		&i.CreatedAt,
	)
	return i, err
}

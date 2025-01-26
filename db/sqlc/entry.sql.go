// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: entry.sql

package db

import (
	"context"
	"time"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountBalance = `-- name: GetAccountBalance :one
SELECT COALESCE(SUM(amount), 0) as total_balance
FROM entries
WHERE account_id = $1
`

func (q *Queries) GetAccountBalance(ctx context.Context, accountID int64) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getAccountBalance, accountID)
	var total_balance interface{}
	err := row.Scan(&total_balance)
	return total_balance, err
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListEntriesParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntries, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
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

const listEntriesByDateRange = `-- name: ListEntriesByDateRange :many
SELECT id, account_id, amount, created_at FROM entries
WHERE 
    account_id = $1 AND
    created_at BETWEEN $2 AND $3
ORDER BY created_at DESC
LIMIT $4
OFFSET $5
`

type ListEntriesByDateRangeParams struct {
	AccountID   int64     `json:"account_id"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedAt_2 time.Time `json:"created_at_2"`
	Limit       int32     `json:"limit"`
	Offset      int32     `json:"offset"`
}

func (q *Queries) ListEntriesByDateRange(ctx context.Context, arg ListEntriesByDateRangeParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntriesByDateRange,
		arg.AccountID,
		arg.CreatedAt,
		arg.CreatedAt_2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
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

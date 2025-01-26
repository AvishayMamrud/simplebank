-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetAccountBalance :one
SELECT COALESCE(SUM(amount), 0) as total_balance
FROM entries
WHERE account_id = $1;

-- name: ListEntriesByDateRange :many
SELECT * FROM entries
WHERE 
    account_id = $1 AND
    created_at BETWEEN $2 AND $3
ORDER BY created_at DESC
LIMIT $4
OFFSET $5;

-- -- name: UpdateEntry :exec
-- UPDATE entries
--   set amount = $2
-- WHERE id = $1;

-- -- name: DeleteEntry :exec
-- DELETE FROM entries
-- WHERE id = $1;
-- name: CreateTransfer :one
INSERT INTO transfers (
  src_account_id,
  dest_account_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE 
    src_account_id = $1 OR
    dest_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: ListTransfersByDateRange :many
SELECT * FROM transfers
WHERE 
    (src_account_id = $1 OR dest_account_id = $2) AND
    created_at BETWEEN $3 AND $4
ORDER BY created_at DESC
LIMIT $5
OFFSET $6;

-- name: GetAccountTransfersSum :one
SELECT 
    COALESCE(SUM(CASE WHEN src_account_id = $1 THEN -amount ELSE amount END), 0) as net_transfer_amount
FROM transfers
WHERE src_account_id = $1 OR dest_account_id = $1;

-- -- name: UpdateTransfer :exec
-- UPDATE transfers
--   set amount = $2
-- WHERE id = $1;

-- -- name: DeleteTransfer :exec
-- DELETE FROM transfers
-- WHERE id = $1;
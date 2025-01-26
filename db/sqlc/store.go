package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

type TransferExecParams struct {
	SrcAccount  int64 `json:"source_account"`
	DestAccount int64 `json:"destination_account"`
	Amount      int64 `json:"amount"`
}

type TransferResultParams struct {
	Transfer    Transfer `json:"transfer"`
	SrcAccount  Account  `json:"source_account"`
	DestAccount Account  `json:"destination_account"`
	SrcEntry    Entry    `json:"source_entry"`
	DestEntry   Entry    `json:"destination_entry"`
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, err2)
		}
		return err
	}

	return tx.Commit()
}

// create transfer create entries for the source account and the dest account
// update balance for both accounts and then finally commit
func (store *Store) TransferTX(ctx context.Context, arg TransferExecParams) (TransferResultParams, error) {
	var result TransferResultParams
	var err error

	err = store.execTx(ctx, func(q *Queries) error {
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			SrcAccountID:  arg.SrcAccount,
			DestAccountID: arg.DestAccount,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.SrcEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: result.Transfer.SrcAccountID,
			Amount:    -result.Transfer.Amount,
		})

		if err != nil {
			return err
		}

		result.DestEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: result.Transfer.DestAccountID,
			Amount:    result.Transfer.Amount,
		})

		if err != nil {
			return err
		}

		// TODO: update accounts balance
		if arg.SrcAccount < arg.DestAccount {
			result.SrcAccount, result.DestAccount, err = AddMoney(ctx, q, arg.SrcAccount, -arg.Amount, arg.DestAccount, arg.Amount)
		} else {
			result.DestAccount, result.SrcAccount, err = AddMoney(ctx, q, arg.DestAccount, arg.Amount, arg.SrcAccount, -arg.Amount)
		}
		return err
	})

	return result, err
}

func AddMoney(
	ctx context.Context,
	q *Queries,
	accID1 int64,
	amount1 int64,
	accID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accID1,
		Amount: amount1,
	})
	if err != nil {
		return // account1, account2, err
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accID2,
		Amount: amount2,
	})
	return // account1, account2, err
}

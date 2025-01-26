package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTX(T *testing.T) {
	// var resultParams TransferResultParams
	// var err error

	_, account1 := createAccountForTest(T)
	_, account2 := createAccountForTest(T)

	store := NewStore(testDB)

	amount := int64(10)

	results := make(chan TransferResultParams)
	errs := make(chan error)

	n := 10
	for i := 0; i < n; i++ {
		go func() {
			resultParams, err := store.TransferTX(context.Background(), TransferExecParams{
				SrcAccount:  account1.ID,
				DestAccount: account2.ID,
				Amount:      amount,
			})

			results <- resultParams
			errs <- err
		}()
	}

	exist := make(map[int]bool)
	for i := 0; i < n; i++ {
		res := <-results
		curr_err := <-errs

		require.NoError(T, curr_err)
		require.NotZero(T, res.Transfer.ID)
		require.NotZero(T, res.Transfer.CreatedAt)
		require.Equal(T, res.Transfer.SrcAccountID, account1.ID)
		require.Equal(T, res.Transfer.DestAccountID, account2.ID)
		require.Equal(T, res.Transfer.Amount, amount)

		require.Equal(T, res.SrcEntry.Amount, -amount)
		require.Equal(T, res.SrcEntry.AccountID, account1.ID)
		require.NotZero(T, res.SrcEntry.CreatedAt)

		require.Equal(T, res.DestEntry.Amount, amount)
		require.Equal(T, res.DestEntry.AccountID, account2.ID)
		require.NotZero(T, res.DestEntry.CreatedAt)

		// check accounts
		require.Equal(T, res.SrcAccount.ID, res.Transfer.SrcAccountID)
		require.Equal(T, res.DestAccount.ID, res.Transfer.DestAccountID)
		require.Equal(T, res.SrcAccount.ID, account1.ID)
		require.Equal(T, res.DestAccount.ID, account2.ID)
		diff1 := account1.Balance - res.SrcAccount.Balance
		diff2 := res.DestAccount.Balance - account2.Balance

		require.Equal(T, diff1, diff2)
		require.Zero(T, diff1%amount)

		k := int(diff1 / amount)
		require.True(T, 0 < k && k <= n)
		require.NotContains(T, exist, k)

		exist[k] = true
	}
	src_acc, err1 := test_queries.GetAccount(context.Background(), account1.ID)
	dst_acc, err2 := test_queries.GetAccount(context.Background(), account2.ID)
	require.NoError(T, err1)
	require.NoError(T, err2)
	require.Equal(T, src_acc.Balance, account1.Balance-int64(n)*amount)
	require.Equal(T, dst_acc.Balance, account2.Balance+int64(n)*amount)
}

func TestDeadlockTransferTX(T *testing.T) {
	// var resultParams TransferResultParams
	// var err error

	_, account1 := createAccountForTest(T)
	_, account2 := createAccountForTest(T)

	store := NewStore(testDB)

	amount := int64(10)

	errs := make(chan error)

	n := 10
	for i := 0; i < n; i++ {
		fromID := account1.ID
		toID := account2.ID
		if i%2 == 0 {
			fromID = account2.ID
			toID = account1.ID
		}

		go func() {
			_, err := store.TransferTX(context.Background(), TransferExecParams{
				SrcAccount:  fromID,
				DestAccount: toID,
				Amount:      amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		curr_err := <-errs
		require.NoError(T, curr_err)
	}
	src_acc, err1 := test_queries.GetAccount(context.Background(), account1.ID)
	dst_acc, err2 := test_queries.GetAccount(context.Background(), account2.ID)
	require.NoError(T, err1)
	require.NoError(T, err2)
	require.Equal(T, src_acc.Balance, account1.Balance)
	require.Equal(T, dst_acc.Balance, account2.Balance)
}

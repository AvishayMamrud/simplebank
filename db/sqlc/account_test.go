package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/mamrud/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createAccountForTest(t *testing.T) (CreateAccountParams, Account) {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := test_queries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	return arg, account
}

func TestCreateAccout(T *testing.T) {
	arg, account := createAccountForTest(T)

	require.NotEmpty(T, account)

	require.Equal(T, arg.Owner, account.Owner)
	require.Equal(T, arg.Balance, account.Balance)
	require.Equal(T, arg.Currency, account.Currency)

	require.NotZero(T, account.ID)
	require.NotZero(T, account.CreatedAt)
}
func TestGetAccount(T *testing.T) {
	_, account := createAccountForTest(T)

	rcv_account, err := test_queries.GetAccount(context.Background(), account.ID)
	require.NoError(T, err)

	require.Equal(T, account.ID, rcv_account.ID)
	require.Equal(T, account.Owner, rcv_account.Owner)
	require.Equal(T, account.Balance, rcv_account.Balance)
	require.Equal(T, account.Currency, rcv_account.Currency)
	require.Equal(T, account.CreatedAt, rcv_account.CreatedAt)

	fmt.Println(rcv_account.CreatedAt)
}
func TestDeleteAccout(T *testing.T) {
	_, account := createAccountForTest(T)

	err := test_queries.DeleteAccount(context.Background(), account.ID)
	require.NoError(T, err)
}
func TestUpdateAccout(T *testing.T) {
	_, account := createAccountForTest(T)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: utils.RandomMoney(),
	}

	_, err := test_queries.UpdateAccount(context.Background(), arg)
	require.NoError(T, err)

	rcv_account, err := test_queries.GetAccount(context.Background(), account.ID)
	require.NoError(T, err)

	require.Equal(T, account.ID, rcv_account.ID)
	require.Equal(T, account.Owner, rcv_account.Owner)
	require.NotEqual(T, account.Balance, rcv_account.Balance)
	require.Equal(T, arg.Balance, rcv_account.Balance)
	require.Equal(T, account.Currency, rcv_account.Currency)
	require.Equal(T, account.CreatedAt, rcv_account.CreatedAt)
}
func TestListAccounts(T *testing.T) {
	for i := 0; i < 10; i++ {
		createAccountForTest(T)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := test_queries.ListAccounts(context.Background(), arg)
	require.NoError(T, err)
	require.Len(T, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(T, account)
	}
}

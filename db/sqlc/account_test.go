package db

import (
	"context"
	"testing"
	"time"

	"github.com/shafi21064/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) Account {
	user := createTestUsers(t)
	arg := CreateAccountParams{
		OwnerName: user.UserName,
		Balance:   util.RandomMoney(),
		Currency:  util.RandomCurrrency(),
	}

	account, err := testQuaries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.OwnerName, account.OwnerName)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createTestAccount(t)

	account2, err := testQuaries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerName, account2.OwnerName)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createTestAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQuaries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerName, account2.OwnerName)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createTestAccount(t)

	err := testQuaries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQuaries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	for range 10 {
		createTestAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQuaries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}

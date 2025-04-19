package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)
	println(">>before: ", account1.Balance, account2.Balance)
	n := 10
	amount := int64(10)

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for range n {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if n%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: pgtype.Int8{Int64: fromAccountID, Valid: true},
				ToAccountId:   pgtype.Int8{Int64: toAccountID, Valid: true},
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	//check result
	for range n {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
		require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID.Int64)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID.Int64)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// Validate database entries exist
		_, err = store.GetTransfer(context.Background(), transfer.FromAccountID)
		require.NoError(t, err)

		_, err = store.GetEntry(context.Background(), fromEntry.AccountID)
		require.NoError(t, err)

		_, err = store.GetEntry(context.Background(), toEntry.AccountID)
		require.NoError(t, err)

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		println(">>tx: ", fromAccount.Balance, toAccount.Balance)

		// check account's balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check final update balance
	updateAccount1, err := testQuaries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQuaries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	println(">>after: ", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

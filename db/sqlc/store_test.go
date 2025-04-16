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

	amount := int64(10)
	result, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountId: pgtype.Int8{Int64: account1.ID, Valid: true},
		ToAccountId:   pgtype.Int8{Int64: account2.ID, Valid: true},
		Amount:        amount,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Verify transfer
	require.NotZero(t, result.Transfer.ID)
	require.Equal(t, account1.ID, result.Transfer.FromAccountID.Int64)
	require.Equal(t, account2.ID, result.Transfer.ToAccountID.Int64)
	require.Equal(t, amount, result.Transfer.Amount)
	require.NotZero(t, result.Transfer.CreatedAt)

	// Verify entries
	require.NotZero(t, result.FromEntry.ID)
	require.Equal(t, account1.ID, result.FromEntry.AccountID.Int64)
	require.Equal(t, -amount, result.FromEntry.Amount)
	require.NotZero(t, result.FromEntry.CreatedAt)

	require.NotZero(t, result.ToEntry.ID)
	require.Equal(t, account2.ID, result.ToEntry.AccountID.Int64)
	require.Equal(t, amount, result.ToEntry.Amount)
	require.NotZero(t, result.ToEntry.CreatedAt)

	// Verify accounts
	require.Equal(t, account1.ID, result.FromAccount.ID)
	require.Equal(t, account2.ID, result.ToAccount.ID)

	// Verify balances
	diff1 := account1.Balance - result.FromAccount.Balance
	diff2 := result.ToAccount.Balance - account2.Balance
	require.Equal(t, diff1, diff2)
	require.True(t, diff1 > 0)
}

// func TestTransferTx(t *testing.T) {
// 	store := NewStore(testDB)

// 	account1 := createTestAccount(t)
// 	account2 := createTestAccount(t)

// 	n := 5 // Number of concurrent transactions
// 	amount := int64(10)

// 	// Channels to collect results
// 	errs := make(chan error, n)
// 	results := make(chan TransferTxResult, n)

// 	// Run concurrent transfer transactions
// 	for i := 0; i < n; i++ {
// 		go func() {
// 			result, err := store.TransferTx(context.Background(), TransferTxParams{
// 				FromAccountId: pgtype.Int8{Int64: account1.ID, Valid: true},
// 				ToAccountId:   pgtype.Int8{Int64: account2.ID, Valid: true},
// 				Amount:        amount,
// 			})
// 			errs <- err
// 			results <- result
// 		}()
// 	}

// 	// Check results
// 	existed := make(map[int]bool)
// 	for i := 0; i < n; i++ {
// 		err := <-errs
// 		require.NoError(t, err)

// 		result := <-results
// 		require.NotEmpty(t, result)

// 		// Check transfer
// 		transfer := result.Transfer
// 		require.NotEmpty(t, transfer)
// 		require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
// 		require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
// 		require.Equal(t, amount, transfer.Amount)
// 		require.NotZero(t, transfer.ID)
// 		require.NotZero(t, transfer.CreatedAt)

// 		// Check entries
// 		fromEntry := result.FromEntry
// 		require.NotEmpty(t, fromEntry)
// 		require.Equal(t, account1.ID, fromEntry.AccountID.Int64)
// 		require.Equal(t, -amount, fromEntry.Amount)
// 		require.NotZero(t, fromEntry.ID)
// 		require.NotZero(t, fromEntry.CreatedAt)

// 		toEntry := result.ToEntry
// 		require.NotEmpty(t, toEntry)
// 		require.Equal(t, account2.ID, toEntry.AccountID.Int64)
// 		require.Equal(t, amount, toEntry.Amount)
// 		require.NotZero(t, toEntry.ID)
// 		require.NotZero(t, toEntry.CreatedAt)

// 		// Check if transfer exists
// 		_, err = store.GetTransfer(context.Background(), transfer.FromAccountID)
// 		require.NoError(t, err)

// 		// Check if entries exist
// 		_, err = store.GetEntry(context.Background(), fromEntry.AccountID)
// 		require.NoError(t, err)

// 		_, err = store.GetEntry(context.Background(), toEntry.AccountID)
// 		require.NoError(t, err)

// 		// Check accounts
// 		fromAccount := result.FromAccount
// 		require.NotEmpty(t, fromAccount)
// 		require.Equal(t, account1.ID, fromAccount.ID)

// 		toAccount := result.ToAccount
// 		require.NotEmpty(t, toAccount)
// 		require.Equal(t, account2.ID, toAccount.ID)

// 		// Check balances
// 		diff1 := account1.Balance - fromAccount.Balance
// 		diff2 := toAccount.Balance - account2.Balance
// 		require.Equal(t, diff1, diff2)
// 		require.True(t, diff1 > 0)
// 		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, ..., n * amount

// 		k := int(diff1 / amount)
// 		require.True(t, k >= 1 && k <= n)
// 		require.NotContains(t, existed, k)
// 		existed[k] = true
// 	}

// 	// Check final updated balances
// 	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
// 	require.NoError(t, err)

// 	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
// 	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
// }

// func TestTransferTx(t *testing.T) {
// 	store := NewStore(testDB)
// 	account1 := createTestAccount(t)
// 	account2 := createTestAccount(t)

// 	n := 5
// 	amount := int64(10)

// 	errs := make(chan error, n)
// 	results := make(chan TransferTxResult, n)

// 	for i := 0; i < n; i++ {
// 		go func() {
// 			result, err := store.TransferTx(context.Background(), TransferTxParams{
// 				FromAccountId: pgtype.Int8{Int64: account1.ID, Valid: true},
// 				ToAccountId:   pgtype.Int8{Int64: account2.ID, Valid: true},
// 				Amount:        amount,
// 			})
// 			errs <- err
// 			results <- result
// 		}()
// 	}

// 	for i := 0; i < n; i++ {
// 		err := <-errs
// 		require.NoError(t, err)

// 		result := <-results
// 		require.NotEmpty(t, result)

// 		// Check transfer
// 		transfer := result.Transfer
// 		require.NotEmpty(t, transfer)
// 		require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
// 		require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
// 		require.Equal(t, amount, transfer.Amount)
// 		require.NotZero(t, transfer.ID)
// 		require.NotZero(t, transfer.CreatedAt)

// 		// Check entries
// 		fromEntry := result.FromEntry
// 		require.NotEmpty(t, fromEntry)
// 		require.Equal(t, account1.ID, fromEntry.AccountID.Int64)
// 		require.Equal(t, -amount, fromEntry.Amount)
// 		require.NotZero(t, fromEntry.ID)
// 		require.NotZero(t, fromEntry.CreatedAt)

// 		toEntry := result.ToEntry
// 		require.NotEmpty(t, toEntry)
// 		require.Equal(t, account2.ID, toEntry.AccountID.Int64)
// 		require.Equal(t, amount, toEntry.Amount)
// 		require.NotZero(t, toEntry.ID)
// 		require.NotZero(t, toEntry.CreatedAt)

// 		// Validate database entries exist
// 		_, err = store.GetTransfer(context.Background(), pgtype.Int8{Int64: transfer.ID, Valid: true})
// 		require.NoError(t, err)

// 		_, err = store.GetEntry(context.Background(), pgtype.Int8{Int64: fromEntry.ID, Valid: true})
// 		require.NoError(t, err)

// 		_, err = store.GetEntry(context.Background(), pgtype.Int8{Int64: toEntry.ID, Valid: true})
// 		require.NoError(t, err)

// 		// TODO: check account balances
// 	}
// }

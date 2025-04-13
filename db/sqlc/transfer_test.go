package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shafi21064/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T) Transfer {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	args := CreateTransferParams{
		FromAccountID: pgtype.Int8{Int64: account1.ID, Valid: true},
		ToAccountID:   pgtype.Int8{Int64: account2.ID, Valid: true},
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQuaries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotNil(t, transfer.ID)
	require.NotNil(t, transfer.FromAccountID)
	require.NotNil(t, transfer.ToAccountID)
	require.NotNil(t, transfer.Amount)

	return transfer
}

func getTransferRecord(t *testing.T, transfer Transfer, searchId *pgtype.Int8) {
	transfer1, err := testQuaries.GetTransfer(context.Background(), *searchId)

	require.NoError(t, err)
	require.NotEmpty(t, transfer1)
	require.Equal(t, transfer.ID, transfer1.ID)
	require.Equal(t, transfer.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer.Amount, transfer1.Amount)
	require.WithinDuration(t, transfer.CreatedAt.Time, transfer1.CreatedAt.Time, time.Second)
}

func TestCreateTransfer(t *testing.T) {
	createTestTransfer(t)
}

func TestGetTransferFromAccountId(t *testing.T) {
	transfer := createTestTransfer(t)
	getTransferRecord(t, transfer, &transfer.FromAccountID)
}
func TestGetTransferToAccountId(t *testing.T) {
	transfer := createTestTransfer(t)
	getTransferRecord(t, transfer, &transfer.ToAccountID)
}

func TestUpdateTransfer(t *testing.T) {
	transfer := createTestTransfer(t)

	arg := UpdateTransferParams{
		FromAccountID: transfer.FromAccountID,
		Amount:        util.RandomMoney(),
	}

	transfer1, err := testQuaries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer1)
	require.Equal(t, transfer.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, arg.Amount, transfer1.Amount)
	require.WithinDuration(t, transfer.CreatedAt.Time, transfer1.CreatedAt.Time, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createTestTransfer(t)

	err := testQuaries.DeleteTransfer(context.Background(), transfer.FromAccountID)
	require.NoError(t, err)
	transfer2, err := testQuaries.GetTransfer(context.Background(), transfer.FromAccountID)
	require.Error(t, err)
	require.Empty(t, transfer2)
}

func TestListTransfer(t *testing.T) {
	for range 10 {
		createTestTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	trandfers, err := testQuaries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trandfers)

	for _, transfer := range trandfers {
		require.NotEmpty(t, transfer)
	}
}

package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shafi21064/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T) Entry {
	account := createTestAccount(t)
	arg := CreateEntryParams{
		AccountID: pgtype.Int8{Int64: account.ID, Valid: true},
		Amount:    util.RandomMoney(),
	}

	entry, err := testQuaries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createTestEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createTestEntry(t)

	entry2, err := testQuaries.GetEntry(context.Background(), entry1.AccountID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	entry1 := createTestEntry(t)

	arg := UpdateEntryParams{
		AccountID: entry1.AccountID,
		Amount:    util.RandomMoney(),
	}

	entry2, err := testQuaries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, time.Second)

}

func TestDeleteEntry(t *testing.T) {
	entry1 := createTestEntry(t)

	err := testQuaries.DeleteEntry(context.Background(), entry1.AccountID)
	require.NoError(t, err)

	entry2, err := testQuaries.GetEntry(context.Background(), entry1.AccountID)
	require.Error(t, err)
	require.Empty(t, entry2)
}

func TestListEntry(t *testing.T) {
	for range 10 {
		createTestEntry(t)
	}
	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQuaries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

}

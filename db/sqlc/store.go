package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shafi21064/simplebank/util"
)

// store provide all functions to execute SQL queries and transactions
type Store struct {
	db *pgxpool.Pool
	*Queries
}

// Newstore creates a new store
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})

	util.CheckError("begin tx:", err)

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	log.Println("Committing transaction") // Debug log
	return tx.Commit(ctx)
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountId pgtype.Int8 `json:"from_account_id"`
	ToAccountId   pgtype.Int8 `json:"to_account_id"`
	Amount        int64       `json:"ammount"`
}

// TransferTxResult is the result of the transaction
type TransferTxResult struct {
	Transfer    `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	FromEntry   Entry   `json:"from_entry"`
	ToEntry     Entry   `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, ans update account balance within a single database transaction
func (strore *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := strore.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		util.CheckError("create transfer error:", err)

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		util.CheckError("", err)

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		util.CheckError("", err)

		// update accounts balance

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountId.Int64,
			Amount: -arg.Amount,
		})

		util.CheckError("", err)

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountId.Int64,
			Amount: arg.Amount,
		})

		util.CheckError("", err)

		if arg.FromAccountId.Int64 < arg.ToAccountId.Int64 {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId.Int64, -arg.Amount, arg.ToAccountId.Int64, arg.Amount)
		} else {
			result.ToAccount, result.ToAccount, err = addMoney(ctx, q, arg.ToAccountId.Int64, arg.Amount, arg.FromAccountId.Int64, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	ammount1 int64,
	accountID2 int64,
	ammount2 int64) (account1 Account,
	account2 Account,
	err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: ammount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: ammount2,
	})
	if err != nil {
		return
	}
	return
}

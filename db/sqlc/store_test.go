package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	fmt.Println(">>before:", acc1.Balance, acc2.Balance)

	n := 5
	amount := int64(1)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func(fromID, toID int64) {
			ctx := context.Background()
			res, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromID,
				ToAccountID:   toID,
				Amount:        amount,
			})

			errs <- err
			results <- res
		}(acc1.ID, acc2.ID)
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		//test transfer
		require.NotEmpty(t, res.Transfer)
		require.Equal(t, acc1.ID, res.Transfer.FromAccountID)
		require.Equal(t, acc2.ID, res.Transfer.ToAccountID)
		require.Equal(t, amount, res.Transfer.Amount)
		require.NotZero(t, res.Transfer.ID)
		require.NotZero(t, res.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), res.Transfer.ID)
		require.NoError(t, err)

		//test outgoing entry
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//test incoming entry
		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)

		//check account balance
		fmt.Println(">>tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
	}

	// check final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	fmt.Println(">>after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, acc1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

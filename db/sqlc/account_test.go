package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tanuj2909/bank/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:   util.RandomOwner(),
		Balance: util.RandomMoney(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Balance, acc.Balance)
	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)

	return acc
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAcc := createRandomAccount(t)

	acc, err := testQueries.GetAccount(context.Background(), newAcc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, newAcc.ID, acc.ID)
	require.Equal(t, newAcc.Balance, acc.Balance)
	require.Equal(t, acc.Owner, newAcc.Owner)
	require.WithinDuration(t, acc.CreatedAt, newAcc.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	newAcc := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      newAcc.ID,
		Balance: util.RandomMoney(),
	}
	acc, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, newAcc.ID, acc.ID)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, acc.Owner, newAcc.Owner)
	require.WithinDuration(t, acc.CreatedAt, newAcc.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	newAcc := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), newAcc.ID)
	require.NoError(t, err)

	acc, err := testQueries.GetAccount(context.Background(), newAcc.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}

	accs, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accs, 5)

	for _, acc := range accs {
		require.NotEmpty(t, acc)
	}
}

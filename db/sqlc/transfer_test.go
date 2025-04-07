package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tanuj2909/bank/util"
)

func createRandomTransfer(t *testing.T) Transfer {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfers(t *testing.T) {
	newTransfer := createRandomTransfer(t)

	arg := GetTransfersParams{
		FromAccountID: newTransfer.FromAccountID,
		ToAccountID:   newTransfer.ToAccountID,
		Limit:         1,
		Offset:        0,
	}
	transfers, err := testQueries.GetTransfers(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, transfers)

	require.Len(t, transfers, 1)
	require.Equal(t, transfers[0].FromAccountID, arg.FromAccountID)
	require.Equal(t, transfers[0].ToAccountID, arg.ToAccountID)
	require.WithinDuration(t, newTransfer.CreatedAt, transfers[0].CreatedAt, time.Second)
}

func TestGetTransfer(t *testing.T) {
	newTransfer := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)

	require.NoError(t, err)

	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, newTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, newTransfer.ToAccountID)
	require.Equal(t, transfer.ID, newTransfer.ID)
	require.WithinDuration(t, transfer.CreatedAt, newTransfer.CreatedAt, time.Second)
}

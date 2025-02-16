//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"

	"backend-trainee-assignment-winter-2025/internal/repo"
	"backend-trainee-assignment-winter-2025/internal/repo/repoerrs"
	"backend-trainee-assignment-winter-2025/tests/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferCoin_Success(t *testing.T) {
	db.SetUp(t)
	defer db.TearDown(t)

	// arrange
	repos := repo.NewRepositories(db.DB)
	ctx := context.Background()

	sender := fixtures.User().Username("sender").Password("1234").V()
	recipient := fixtures.User().Username("recipient").Password("1234").V()

	senderResp, err := repos.User.Create(ctx, sender)
	require.NoError(t, err)
	recipientResp, err := repos.User.Create(ctx, recipient)
	require.NoError(t, err)

	transfer := fixtures.Transaction().FromUser(senderResp.Id).ToUser(recipientResp.Id).Amount(30).V()

	// act
	err = repos.Transaction.Transfer(ctx, transfer, repos.User.GetById)

	// assert
	require.NoError(t, err)
}

func TestTransferCoin_InsufficientFunds(t *testing.T) {
	db.SetUp(t)
	defer db.TearDown(t)

	// arrange
	repos := repo.NewRepositories(db.DB)
	ctx := context.Background()

	sender := fixtures.User().Username("sender").Password("1234").V()
	recipient := fixtures.User().Username("recipient").Password("1234").V()

	senderResp, err := repos.User.Create(ctx, sender)
	require.NoError(t, err)
	recipientResp, err := repos.User.Create(ctx, recipient)
	require.NoError(t, err)

	transfer := fixtures.Transaction().FromUser(senderResp.Id).ToUser(recipientResp.Id).Amount(1050).V()

	// act
	err = repos.Transaction.Transfer(ctx, transfer, repos.User.GetById)

	// assert
	require.Error(t, err)
	assert.ErrorIs(t, err, repoerrs.ErrLowBalance)
}

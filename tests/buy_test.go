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

func TestBuyMerch(t *testing.T) {
	var (
		ctx context.Context
	)
	t.Run(
		"success", func(t *testing.T) {
			db.SetUp(t)
			defer db.TearDown(t)
			// arrange
			repos := repo.NewRepositories(db.DB)
			ctx = context.Background()

			user := fixtures.User().Valid().V()
			//fmt.Println(fmt.Sprintf("user: %v", user))

			// добавляем пользователя
			resp, err := repos.User.Create(ctx, user)

			// act
			inventory := fixtures.Inventory().CustomerId(resp.ID).Type("cup").V()
			res, err := repos.Inventory.Add(ctx, inventory)

			// assert
			require.NoError(t, err)
			assert.Equal(t, inventory.CustomerId, res.CustomerID)
			assert.Equal(t, inventory.Type, res.Type)
			assert.Equal(t, 1, res.Quantity)
		},
	)
}

func TestBuyMerch_InsufficientFunds(t *testing.T) {
	var (
		ctx context.Context
	)
	t.Run(
		"insufficient funds", func(t *testing.T) {
			db.SetUp(t)
			defer db.TearDown(t)
			// arrange
			repos := repo.NewRepositories(db.DB)
			ctx = context.Background()

			user := fixtures.User().Valid().V()
			resp, err := repos.User.Create(ctx, user)
			require.NoError(t, err)

			// act
			inventory := fixtures.Inventory().CustomerId(resp.ID).Type("pink-hoody").V()
			_, err = repos.Inventory.Add(ctx, inventory)
			_, err = repos.Inventory.Add(ctx, inventory)
			_, err = repos.Inventory.Add(ctx, inventory)

			// assert
			require.Error(t, err)
			assert.ErrorIs(t, err, repoerrs.ErrLowBalance)
		},
	)
}

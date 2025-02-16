//go:build integration
// +build integration

package tests

import (
	"context"
	"fmt"
	"testing"

	"backend-trainee-assignment-winter-2025/internal/repo"
	"backend-trainee-assignment-winter-2025/tests/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
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

			// act
			user := fixtures.User().Valid().V()

			fmt.Println(fmt.Sprintf("user: %v", user))
			resp, err := repos.User.Create(ctx, user)

			// assert
			require.NoError(t, err)
			assert.Equal(t, resp.Username, user.Username)
		},
	)
}

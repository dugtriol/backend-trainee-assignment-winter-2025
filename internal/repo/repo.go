package repo

import (
	"backend-trainee-assignment-winter-2025/internal/repo/pgdb"
	"backend-trainee-assignment-winter-2025/pkg/postgres"
)

type User interface {
	//Create(ctx context.Context, input entity.User) (string, error)
	//GetById(ctx context.Context, id string) (entity.User, error)
	//GetByUsername(ctx context.Context, username string) (entity.User, error)
}

type Inventory interface {
}

type Transaction interface {
}

type Repositories struct {
	User
	Inventory
	Transaction
}

func NewRepositories(db *postgres.Database) *Repositories {
	return &Repositories{
		User:        pgdb.NewUserRepository(db),
		Inventory:   pgdb.NewInventoryRepository(db),
		Transaction: pgdb.NewTransactionRepository(db),
	}
}

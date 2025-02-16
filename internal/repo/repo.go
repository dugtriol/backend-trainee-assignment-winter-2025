package repo

import (
	"context"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo/pgdb"
	"backend-trainee-assignment-winter-2025/pkg/postgres"
)

type User interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetById(ctx context.Context, id string) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
}

type Inventory interface {
	Add(ctx context.Context, inventory entity.Inventory) (entity.Inventory, error)
	GetByUserID(ctx context.Context, userId string) ([]entity.Inventory, error)
}

type Transaction interface {
	Transfer(
		ctx context.Context, input entity.Transaction,
		isExist func(ctx context.Context, id string) (entity.User, error),
	) error
	GetByUserID(ctx context.Context, userId string) ([]entity.Transaction, error)
}

type Merch interface {
	GetById(ctx context.Context, id string) (entity.Merch, error)
	GetByName(ctx context.Context, name string) (entity.Merch, error)
}

type Repositories struct {
	User
	Inventory
	Transaction
	Merch
}

func NewRepositories(db *postgres.Database) *Repositories {
	return &Repositories{
		User:        pgdb.NewUserRepository(db),
		Inventory:   pgdb.NewInventoryRepository(db),
		Transaction: pgdb.NewTransactionRepository(db),
		Merch:       pgdb.NewMerchRepository(db),
	}
}

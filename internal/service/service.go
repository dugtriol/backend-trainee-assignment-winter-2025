package service

import (
	"context"
	"log/slog"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo"
)

type AuthInput struct {
	Username string
	Password string
}

type User interface {
	Auth(ctx context.Context, log *slog.Logger, input AuthInput) (string, error)
	GetById(ctx context.Context, log *slog.Logger, id string) (entity.User, error)
}

type Inventory interface {
	GetItem(ctx context.Context, log *slog.Logger, userId, item string) error
}

type TransactionInput struct {
	FromUserId string
	ToUserId   string
	Amount     int
}

type Transaction interface {
	Transfer(
		ctx context.Context, log *slog.Logger, input TransactionInput,
	) error
}

type Services struct {
	User
	Inventory
	Transaction
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(dep ServicesDependencies) *Services {
	return &Services{
		User:        NewUserService(dep.Repos.User),
		Inventory:   NewInventoryService(dep.Repos.Inventory),
		Transaction: NewTransactionService(dep.Repos.Transaction, dep.Repos.User),
	}
}

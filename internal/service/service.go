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
	GetByUserId(ctx context.Context, log *slog.Logger, userId string) (
		[]entity.Inventory, error,
	)
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
	GetByUserId(
		ctx context.Context, log *slog.Logger, userId string,
	) ([]entity.Transaction, error)
}

type Info interface {
	Get(
		ctx context.Context, log *slog.Logger, userId string,
		inventories []entity.Inventory,
		transactions []entity.Transaction,
	) ([]InfoInventory, CoinHistory, error)
}
type Services struct {
	User
	Inventory
	Transaction
	Info
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(dep ServicesDependencies) *Services {
	return &Services{
		User:        NewUserService(dep.Repos.User),
		Inventory:   NewInventoryService(dep.Repos.Inventory),
		Transaction: NewTransactionService(dep.Repos.Transaction, dep.Repos.User),
		Info:        NewInfoService(),
	}
}

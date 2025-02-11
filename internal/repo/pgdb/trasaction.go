package pgdb

import "backend-trainee-assignment-winter-2025/pkg/postgres"

const (
	transactionTable = "transactions"
)

type TransactionRepository struct {
	*postgres.Database
}

func NewTransactionRepository(db *postgres.Database) *TransactionRepository {
	return &TransactionRepository{db}
}

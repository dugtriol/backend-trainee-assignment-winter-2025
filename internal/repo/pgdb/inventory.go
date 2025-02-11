package pgdb

import "backend-trainee-assignment-winter-2025/pkg/postgres"

const (
	inventoryTable = "inventories"
)

type InventoryRepository struct {
	*postgres.Database
}

func NewInventoryRepository(db *postgres.Database) *InventoryRepository {
	return &InventoryRepository{db}
}

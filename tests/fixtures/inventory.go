//go:build integration
// +build integration

package fixtures

import "backend-trainee-assignment-winter-2025/internal/entity"

type InventoryBuilder struct {
	instance *entity.Inventory
}

func Inventory() *InventoryBuilder {
	return &InventoryBuilder{instance: &entity.Inventory{}}
}

func (b *InventoryBuilder) CustomerId(value string) *InventoryBuilder {
	b.instance.CustomerId = value
	return b
}

func (b *InventoryBuilder) Type(value string) *InventoryBuilder {
	b.instance.Type = value
	return b
}

func (b *InventoryBuilder) P() *entity.Inventory {
	return b.instance
}

func (b *InventoryBuilder) V() entity.Inventory {
	return *b.instance
}

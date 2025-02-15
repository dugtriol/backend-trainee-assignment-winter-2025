package service

import (
	"context"
	"fmt"
	"log/slog"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo"
)

type InventoryService struct {
	inventoryRepo repo.Inventory
}

func NewInventoryService(inventoryRepo repo.Inventory) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

func (s *InventoryService) сheckParamItem(item string) bool {
	if item == "" || len(item) == 0 {
		return false
	}
	return true
}

func (s *InventoryService) GetItem(ctx context.Context, log *slog.Logger, userId, item string) error {
	// проверка корректности поля
	if !s.сheckParamItem(item) {
		log.Error(fmt.Sprintf("Service - InventoryService - GetItem - сheckParamItem: %s", item))
		return ErrInvalidMerchType
	}
	// проверка корректности типа мерча
	// проверка баланса пользователя
	// покупка мерча
	// изменение баланса у пользователя
	if err := s.inventoryRepo.Add(ctx, entity.Inventory{CustomerId: userId, Type: item}); err != nil {
		return err
	}
	return nil
}

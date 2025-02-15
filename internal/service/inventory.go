package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo"
	"backend-trainee-assignment-winter-2025/internal/repo/repoerrs"
)

type InventoryService struct {
	inventoryRepo repo.Inventory
}

func NewInventoryService(inventoryRepo repo.Inventory) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

func (s *InventoryService) checkParamItem(item string) bool {
	if item == "" || len(item) == 0 {
		return false
	}
	return true
}

func (s *InventoryService) GetItem(ctx context.Context, log *slog.Logger, userId, item string) error {
	if !s.checkParamItem(item) {
		log.Error(fmt.Sprintf("Service - InventoryService - GetItem - сheckParamItem: %s", item))
		return ErrInvalidMerchType
	}
	// проверка корректности типа мерча
	// проверка баланса пользователя
	// покупка мерча
	// изменение баланса у пользователя
	if err := s.inventoryRepo.Add(ctx, entity.Inventory{CustomerId: userId, Type: item}); err != nil {
		if errors.Is(err, repoerrs.ErrLowBalance) {
			return ErrLowBalance
		}
		return err
	}
	return nil
}

func (s *InventoryService) GetByUserId(ctx context.Context, log *slog.Logger, userId string) (
	[]entity.Inventory, error,
) {
	var err error
	var inventories []entity.Inventory
	if inventories, err = s.inventoryRepo.GetByUserID(ctx, userId); err != nil {
		log.Error(fmt.Sprintf("Service - InventoryService - GetByUserId: %v", err))
		return []entity.Inventory{}, err
	}

	if len(inventories) == 0 {
		return []entity.Inventory{}, ErrNotFound
	}
	return inventories, nil
}

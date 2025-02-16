package service_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockInventoryRepo struct {
	mock.Mock
}

func (m *MockInventoryRepo) Add(ctx context.Context, inventory entity.Inventory) (entity.Inventory, error) {
	args := m.Called(ctx, inventory)
	return args.Get(0).(entity.Inventory), args.Error(1)
}

func (m *MockInventoryRepo) GetByUserID(ctx context.Context, userId string) ([]entity.Inventory, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]entity.Inventory), args.Error(1)
}

func TestGetItem_Success(t *testing.T) {
	mockRepo := new(MockInventoryRepo)
	serviceInventory := service.NewInventoryService(mockRepo)
	ctx := context.Background()
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	mockRepo.On("Add", mock.Anything, entity.Inventory{CustomerID: "user123", Type: "hat"}).Return(
		entity.Inventory{
			CustomerID: "user123", Type: "hat", Quantity: 1,
		}, nil,
	)

	err := serviceInventory.GetItem(ctx, log, "user123", "hat")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetItem_LowBalance(t *testing.T) {
	mockRepo := new(MockInventoryRepo)
	serviceInventory := service.NewInventoryService(mockRepo)
	ctx := context.Background()
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	mockRepo.On(
		"Add",
		mock.Anything,
		entity.Inventory{CustomerID: "user123", Type: "hat"},
	).Return(entity.Inventory{}, service.ErrLowBalance)

	err := serviceInventory.GetItem(ctx, log, "user123", "hat")

	assert.ErrorIs(t, err, service.ErrLowBalance)
}

func TestGetByUserId_Success(t *testing.T) {
	mockRepo := new(MockInventoryRepo)
	serviceInventory := service.NewInventoryService(mockRepo)
	ctx := context.Background()
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	mockRepo.On("GetByUserID", mock.Anything, "user123").Return(
		[]entity.Inventory{
			{
				CustomerID: "user123", Type: "hat", Quantity: 1,
			},
		}, nil,
	)

	items, err := serviceInventory.GetByUserID(ctx, log, "user123")

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "hat", items[0].Type)
}

func TestGetByUserId_NotFound(t *testing.T) {
	mockRepo := new(MockInventoryRepo)
	serviceInventory := service.NewInventoryService(mockRepo)
	ctx := context.Background()
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	mockRepo.On("GetByUserID", mock.Anything, "user123").Return([]entity.Inventory{}, nil)

	items, err := serviceInventory.GetByUserID(ctx, log, "user123")

	assert.ErrorIs(t, err, service.ErrNotFound)
	assert.Empty(t, items)
}

package service_test

import (
	"context"
	"log/slog"
	"reflect"
	"testing"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/service"
)

func TestInfoService_Get(t *testing.T) {
	s := service.NewInfoService()
	ctx := context.Background()
	logger := slog.Default()
	userID := "user1"

	inventories := []entity.Inventory{
		{Type: "gold", Quantity: 10},
		{Type: "silver", Quantity: 5},
	}

	transactions := []entity.Transaction{
		{FromUser: "user1", ToUser: "user2", Amount: 50},
		{FromUser: "user3", ToUser: "user1", Amount: 30},
	}

	expectedInventory := []service.InfoInventory{
		{Type: "gold", Quantity: 10},
		{Type: "silver", Quantity: 5},
	}

	expectedCoinHistory := service.CoinHistory{
		Received: []service.TransactionReceived{{FromUser: "user3", Amount: 30}},
		Sent:     []service.TransactionSent{{ToUser: "user2", Amount: 50}},
	}

	gotInventory, gotCoinHistory, err := s.Get(ctx, logger, userID, inventories, transactions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(gotInventory, expectedInventory) {
		t.Errorf("expected inventory: %v, got: %v", expectedInventory, gotInventory)
	}

	if !reflect.DeepEqual(gotCoinHistory, expectedCoinHistory) {
		t.Errorf("expected coin history: %v, got: %v", expectedCoinHistory, gotCoinHistory)
	}
}

func TestInfoService_Get_EmptyData(t *testing.T) {
	s := service.NewInfoService()
	ctx := context.Background()
	logger := slog.Default()
	userID := "user1"

	gotInventory, gotCoinHistory, err := s.Get(ctx, logger, userID, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(gotInventory) != 0 {
		t.Errorf("expected empty inventory, got: %v", gotInventory)
	}

	if len(gotCoinHistory.Received) != 0 || len(gotCoinHistory.Sent) != 0 {
		t.Errorf("expected empty coin history, got: %v", gotCoinHistory)
	}
}

package service

import (
	"context"
	"log/slog"

	"backend-trainee-assignment-winter-2025/internal/entity"
)

type InfoService struct {
	//transactionRepo repo.Transaction
	//inventoryRepo   repo.Inventory
	//userRepo        repo.User
}

func NewInfoService() *InfoService {
	return &InfoService{}
}

func (s *InfoService) Get(
	ctx context.Context, log *slog.Logger, userID string,
	inventories []entity.Inventory,
	transactions []entity.Transaction,
) ([]InfoInventory, CoinHistory, error) {
	var err error
	var infoInventory []InfoInventory
	var coinHistory CoinHistory

	if len(inventories) != 0 {
		if infoInventory = s.getInventoryArray(inventories); err != nil {
			return []InfoInventory{}, CoinHistory{}, err
		}
	}

	if len(transactions) != 0 {
		if coinHistory = s.getTransactionArray(transactions, userID); err != nil {
			return []InfoInventory{}, CoinHistory{}, err
		}
	}
	return infoInventory, coinHistory, nil
}

func (s *InfoService) getInventoryArray(inventories []entity.Inventory) []InfoInventory {
	infoInventories := make([]InfoInventory, len(inventories))
	for i, inv := range inventories {
		infoInventories[i] = InfoInventory{
			Type:     inv.Type,
			Quantity: inv.Quantity,
		}
	}

	return infoInventories
}

type InfoInventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []TransactionReceived `json:"received"`
	Sent     []TransactionSent     `json:"sent"`
}

type TransactionReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type TransactionSent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (s *InfoService) getTransactionArray(transactions []entity.Transaction, userID string) CoinHistory {
	trReceived := make([]TransactionReceived, 0)
	trSent := make([]TransactionSent, 0)

	coinHistoryArr := CoinHistory{
		Received: trReceived,
		Sent:     trSent,
	}

	for _, inv := range transactions {
		if inv.FromUser == userID {
			coinHistoryArr.Sent = append(
				coinHistoryArr.Sent, TransactionSent{
					ToUser: inv.ToUser,
					Amount: inv.Amount,
				},
			)
		} else {
			coinHistoryArr.Received = append(
				coinHistoryArr.Received, TransactionReceived{
					FromUser: inv.FromUser,
					Amount:   inv.Amount,
				},
			)
		}
	}

	return coinHistoryArr
}

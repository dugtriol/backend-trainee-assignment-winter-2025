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

type TransactionService struct {
	transactionRepo repo.Transaction
	userRepo        repo.User
}

func NewTransactionService(transactionRepo repo.Transaction, userRepo repo.User) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo, userRepo: userRepo}
}

func (s *TransactionService) Transfer(
	ctx context.Context, log *slog.Logger, input TransactionInput,
) error {
	if input.FromUserID == input.ToUserID {
		return ErrSimilarID
	}
	transaction := entity.Transaction{
		FromUser: input.FromUserID,
		ToUser:   input.ToUserID,
		Amount:   input.Amount,
	}
	if err := s.transactionRepo.Transfer(ctx, transaction, s.userRepo.GetByID); err != nil {
		log.Error(fmt.Sprintf("Service - TransactionService - Transfer - transactionRepo.Transfer: %v", err))
		if errors.Is(err, repoerrs.ErrLowBalance) {
			return ErrLowBalance
		}
		return err
	}
	return nil
}

func (s *TransactionService) GetByUserID(
	ctx context.Context, log *slog.Logger, userID string,
) ([]entity.Transaction, error) {
	var err error
	var transactions []entity.Transaction
	if transactions, err = s.transactionRepo.GetByUserID(ctx, userID); err != nil {
		log.Error(fmt.Sprintf("Service - TransactionService - GetByUserID: %v", err))
		return []entity.Transaction{}, err
	}

	if len(transactions) == 0 {
		return []entity.Transaction{}, ErrNotFound
	}
	return transactions, nil
}

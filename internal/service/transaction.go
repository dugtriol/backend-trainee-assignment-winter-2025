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
	if input.FromUserId == input.ToUserId {
		return ErrSimilarId
	}
	transaction := entity.Transaction{
		FromUser: input.FromUserId,
		ToUser:   input.ToUserId,
		Amount:   input.Amount,
	}
	if err := s.transactionRepo.Transfer(ctx, transaction, s.userRepo.GetById); err != nil {
		log.Error(fmt.Sprintf("Service - TransactionService - Transfer - transactionRepo.Transfer: %v", err))
		if errors.Is(err, repoerrs.ErrLowBalance) {
			return ErrLowBalance
		}
		return err
	}
	return nil
}

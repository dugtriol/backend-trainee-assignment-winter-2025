package service_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo/repoerrs"
	"backend-trainee-assignment-winter-2025/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) Transfer(
	ctx context.Context, tx entity.Transaction, isExist func(context.Context, string) (entity.User, error),
) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepo) GetByUserID(ctx context.Context, userId string) ([]entity.Transaction, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user entity.User) (entity.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.User), args.Error(1)
}

func TestTransactionService_Transfer_Success(t *testing.T) {
	transactionRepo := new(MockTransactionRepo)
	userRepo := new(MockUserRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	input := service.TransactionInput{
		FromUserID: "user1",
		ToUserID:   "user2",
		Amount:     100,
	}

	transactionRepo.On("Transfer", mock.Anything, mock.Anything).Return(nil)

	err := transactionService.Transfer(context.Background(), log, input)

	assert.NoError(t, err)
	transactionRepo.AssertExpectations(t)
}

func TestTransactionService_Transfer_SameUser(t *testing.T) {
	transactionService := service.NewTransactionService(nil, nil)
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	input := service.TransactionInput{
		FromUserID: "user1",
		ToUserID:   "user1",
		Amount:     100,
	}

	err := transactionService.Transfer(context.Background(), log, input)

	assert.ErrorIs(t, err, service.ErrSimilarID)
}

func TestTransactionService_Transfer_LowBalance(t *testing.T) {
	transactionRepo := new(MockTransactionRepo)
	userRepo := new(MockUserRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	input := service.TransactionInput{
		FromUserID: "user1",
		ToUserID:   "user2",
		Amount:     1000,
	}

	transactionRepo.On("Transfer", mock.Anything, mock.Anything).Return(repoerrs.ErrLowBalance)

	err := transactionService.Transfer(context.Background(), log, input)

	assert.ErrorIs(t, err, service.ErrLowBalance)
}

func TestTransactionService_GetByUserId_Success(t *testing.T) {
	transactionRepo := new(MockTransactionRepo)
	transactionService := service.NewTransactionService(transactionRepo, nil)
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	transactions := []entity.Transaction{{ID: "txn1", FromUser: "user1", ToUser: "user2", Amount: 100}}
	transactionRepo.On("GetByUserID", mock.Anything, "user1").Return(transactions, nil)

	result, err := transactionService.GetByUserID(context.Background(), log, "user1")

	assert.NoError(t, err)
	assert.Equal(t, transactions, result)
}

func TestTransactionService_GetByUserId_NotFound(t *testing.T) {
	transactionRepo := new(MockTransactionRepo)
	transactionService := service.NewTransactionService(transactionRepo, nil)
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	transactionRepo.On("GetByUserID", mock.Anything, "user1").Return([]entity.Transaction{}, nil)

	_, err := transactionService.GetByUserID(context.Background(), log, "user1")

	assert.ErrorIs(t, err, service.ErrNotFound)
}

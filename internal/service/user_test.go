package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/pkg/hasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestUserService_Auth_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	ctx := context.Background()
	hashedPassword, _ := hasher.HashPassword("password123")
	user := entity.User{ID: "user123", Username: "testuser", Password: hashedPassword}

	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(user, nil)

	tokenString, err := userService.Auth(
		ctx, log, AuthInput{
			Username: "testuser",
			Password: "password123",
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Auth_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	ctx := context.Background()
	hashedPassword, _ := hasher.HashPassword("correct_password")
	user := entity.User{ID: "user123", Username: "testuser", Password: hashedPassword}

	mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(user, nil)

	_, err := userService.Auth(
		ctx, log, AuthInput{
			Username: "testuser",
			Password: "wrong_password",
		},
	)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassword, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	ctx := context.Background()
	input := AuthInput{Username: "newuser", Password: "password123"}
	hashedPassword, _ := hasher.HashPassword(input.Password)
	createdUser := entity.User{ID: "user456", Username: input.Username, Password: hashedPassword}

	mockRepo.On("GetByUsername", mock.Anything, input.Username).Return(entity.User{}, ErrUserNotFound)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(createdUser, nil)

	tokenString, err := userService.Auth(ctx, log, input)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetById_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)
	ctx := context.Background()
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	user := entity.User{ID: "user789", Username: "someuser", Password: "hashedpass"}
	mockRepo.On("GetByID", ctx, "user789").Return(user, nil)

	result, err := userService.GetByID(ctx, log, "user789")

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetById_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)
	ctx := context.Background()
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	mockRepo.On("GetByID", ctx, "nonexistent").Return(entity.User{}, ErrUserNotFound)

	_, err := userService.GetByID(ctx, log, "nonexistent")

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

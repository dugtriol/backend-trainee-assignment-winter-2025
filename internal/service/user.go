package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo"
	"backend-trainee-assignment-winter-2025/internal/repo/repoerrs"
	"backend-trainee-assignment-winter-2025/pkg/hasher"
	"backend-trainee-assignment-winter-2025/pkg/token"
)

type UserService struct {
	userRepo repo.User
}

func NewUserService(userRepo repo.User) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Auth(ctx context.Context, log *slog.Logger, input AuthInput) (string, error) {
	log.Info(fmt.Sprintf("Service - UserService - Auth"))
	var err error
	var tokenString string
	var output entity.User
	output, err = s.isExist(ctx, log, input)

	if errors.Is(err, ErrInvalidPassword) {
		return "", err
	}
	if err != nil {
		if tokenString, err = token.Create(output.Id); err != nil {
			return "", err
		}
		return tokenString, err
	}

	if tokenString, err = s.register(ctx, log, input); err != nil {
		return "", err
	}
	return tokenString, err
}

func (s *UserService) GetById(ctx context.Context, log *slog.Logger, id string) (entity.User, error) {
	var err error
	log.Info(fmt.Sprintf("Service - UserService - GetById"))
	output, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - GetById - GetById: %v", err))
		return entity.User{}, ErrUserNotFound
	}
	return output, err
}

func (s *UserService) register(ctx context.Context, log *slog.Logger, input AuthInput) (string, error) {
	var err error
	var tokenString string
	log.Info(fmt.Sprintf("Service - UserService - register"))
	password, err := hasher.HashPassword(input.Password)
	if err != nil {
		return "", ErrCannotHashPassword
	}
	user := entity.User{
		Username: input.Username,
		Password: password,
	}
	var output entity.User
	output, err = s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return "", ErrUserAlreadyExists
		}
		log.Error("UserService.Register - c.userRepo.Register: %v", err)
		return "", ErrCannotCreateUser
	}

	if tokenString, err = token.Create(output.Id); err != nil {
		return "", err
	}
	return tokenString, err
}

func (s *UserService) isExist(ctx context.Context, log *slog.Logger, input AuthInput) (entity.User, error) {
	var err error
	log.Info(fmt.Sprintf("Service - UserService - isExist"))
	output, err := s.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		log.Error(fmt.Sprintf("Service - UserService - isExist - GetByUsername: %v", err))
		return entity.User{}, ErrUserNotFound
	}

	if err = hasher.CheckPassword(input.Password, output.Password); err != nil {
		log.Error(fmt.Sprintf("Service - UserService - isExist - CheckPassword: %v", err))
		return entity.User{}, ErrInvalidPassword
	}

	return output, err
}

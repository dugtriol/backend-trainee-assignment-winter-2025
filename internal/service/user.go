package service

import "backend-trainee-assignment-winter-2025/internal/repo"

type UserService struct {
	userRepo repo.User
}

func NewUserService(userRepo repo.User) *UserService {
	return &UserService{userRepo: userRepo}
}

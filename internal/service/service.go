package service

import (
	"context"
	"log/slog"

	"backend-trainee-assignment-winter-2025/internal/repo"
)

type AuthInput struct {
	Username string
	Password string
}

type User interface {
	Auth(ctx context.Context, log *slog.Logger, input AuthInput) (string, error)
}

type Services struct {
	User
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(dep ServicesDependencies) *Services {
	return &Services{
		User: NewUserService(dep.Repos.User),
	}
}

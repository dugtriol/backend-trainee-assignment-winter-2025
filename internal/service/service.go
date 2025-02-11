package service

import "backend-trainee-assignment-winter-2025/internal/repo"

type User interface{}

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

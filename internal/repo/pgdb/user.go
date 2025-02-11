package pgdb

import "backend-trainee-assignment-winter-2025/pkg/postgres"

const (
	userTable = "users"
)

type UserRepository struct {
	*postgres.Database
}

func NewUserRepository(db *postgres.Database) *UserRepository {
	return &UserRepository{db}
}

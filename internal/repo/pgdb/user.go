package pgdb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo/repoerrs"
	"backend-trainee-assignment-winter-2025/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

const (
	userTable = "users"
)

type UserRepository struct {
	*postgres.Database
}

func NewUserRepository(db *postgres.Database) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) Create(ctx context.Context, user entity.User) error {
	var err error
	sql, args, err := u.Builder.Insert(userTable).Columns("username", "password").Values(
		user.Username,
		user.Password,
	).ToSql()

	log.Println(sql)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - u.Builder.Insert: %v", err)
	}

	if _, err = u.Cluster.Exec(ctx, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerrs.ErrNotFound
		}
		return fmt.Errorf("UserRepo - Create: %v", err)
	}

	return nil
}

func (u *UserRepository) GetById(ctx context.Context, id string) (entity.User, error) {
	return u.getByField(ctx, "id", id)
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	return u.getByField(ctx, "username", username)
}

func (u *UserRepository) getByField(ctx context.Context, field, value string) (entity.User, error) {
	var err error
	sql, args, _ := u.Builder.
		Select("*").
		From(userTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		ToSql()
	log.Printf("UserRepo - GetByField - sql %s args %s \n", sql, args)

	var output entity.User
	err = u.Cluster.QueryRow(ctx, sql, args...).Scan(
		&output.Id,
		&output.Username,
		&output.Password,
		&output.Amount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo - GetByField %s - r.Cluster.QueryRow: %v", field, err)
	}
	return output, nil
}

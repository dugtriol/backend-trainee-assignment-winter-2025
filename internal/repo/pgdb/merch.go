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
	merchTable = "merch"
)

type MerchRepository struct {
	*postgres.Database
}

func NewMerchRepository(db *postgres.Database) *MerchRepository {
	return &MerchRepository{db}
}

func (u *MerchRepository) GetById(ctx context.Context, id string) (entity.Merch, error) {
	return u.getByField(ctx, "id", id)
}

func (u *MerchRepository) GetByName(ctx context.Context, name string) (entity.Merch, error) {
	return u.getByField(ctx, "name", name)
}

func (u *MerchRepository) getByField(ctx context.Context, field, value string) (entity.Merch, error) {
	var err error
	sql, args, _ := u.Builder.
		Select("*").
		From(merchTable).
		Where(fmt.Sprintf("%v = ?", field), value).
		ToSql()
	log.Printf("MerchRepository - getByField - sql %s args %s \n", sql, args)

	var output entity.Merch
	err = u.Cluster.QueryRow(ctx, sql, args...).Scan(
		&output.Id,
		&output.Name,
		&output.Price,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Merch{}, repoerrs.ErrNotFound
		}
		return entity.Merch{}, fmt.Errorf("MerchRepository - getByField %s - r.Cluster.QueryRow: %v", field, err)
	}
	return output, nil
}

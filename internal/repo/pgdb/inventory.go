package pgdb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo/repoerrs"
	"backend-trainee-assignment-winter-2025/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	inventoryTable = "inventories"
)

type InventoryRepository struct {
	*postgres.Database
}

func NewInventoryRepository(db *postgres.Database) *InventoryRepository {
	return &InventoryRepository{db}
}

func (u *InventoryRepository) Add(ctx context.Context, inventory entity.Inventory) error {
	var err error
	tx, err := u.Cluster.Begin(ctx)
	if err != nil {
		return fmt.Errorf("InventoryRepository - starting transaction: %w", err)
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit(ctx)
		} else {
			e = tx.Rollback(ctx)
		}

		if err == nil && e != nil {
			err = fmt.Errorf("InventoryRepository - finishing transaction: %w", e)
		}
	}()

	// Вычитаем баланс пользователя
	if err = u.buyMerch(ctx, tx, inventory.CustomerId, inventory.Type); err != nil {
		return err
	}

	// Обновляем инвентарь
	if err = u.updateInventory(ctx, tx, inventory.CustomerId, inventory.Type); err != nil {
		return err
	}

	return nil
}

func (u *InventoryRepository) buyMerch(
	ctx context.Context, tx pgx.Tx, customerId string, merchType string,
) error {
	var err error
	// Получаем цену товара
	queryPrice, argsPrice, err := u.Builder.
		Select("price").
		From(merchTable).
		Where(squirrel.Eq{"name": merchType}).
		ToSql()
	log.Printf("InventoryRepository - getByField - sql %s args %s \n", queryPrice, argsPrice)
	if err != nil {
		return fmt.Errorf("InventoryRepository - buyMerch - building price query: %w", err)
	}

	var price int
	if err = tx.QueryRow(ctx, queryPrice, argsPrice...).Scan(&price); err != nil {
		return fmt.Errorf("InventoryRepository - buyMerch - getting merch price: %w", err)
	}

	// Обновляем баланс пользователя
	queryAmount, argsAmount, err := u.Builder.
		Update(userTable).
		Set("amount", squirrel.Expr("amount - ?", price)).
		Where(
			squirrel.And{
				squirrel.Eq{"id": customerId},
				squirrel.GtOrEq{"amount": price}, // Проверяем, хватает ли баланса
			},
		).
		Suffix("RETURNING amount").
		ToSql()

	log.Printf("InventoryRepository - getByField - sql %s args %s \n", queryAmount, argsAmount)
	if err != nil {
		return fmt.Errorf("buyMerch - building amount update query: %w", err)
	}

	var updatedAmount int
	if err = tx.QueryRow(ctx, queryAmount, argsAmount...).Scan(&updatedAmount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerrs.ErrLowBalance
		}
		return fmt.Errorf("buyMerch - updating user amount: %w", err)
	}

	return nil
}

// Добавляет товар в инвентарь или увеличивает количество
func (u *InventoryRepository) updateInventory(
	ctx context.Context, tx pgx.Tx, customerId string, merchType string,
) error {
	queryInventory, argsInventory, err := u.Builder.Insert(inventoryTable).
		Columns("customer_id", "type", "quantity").
		Values(customerId, merchType, 1).
		Suffix("ON CONFLICT (customer_id, type) DO UPDATE SET quantity = inventories.quantity + 1").
		ToSql()
	log.Printf("InventoryRepository - getByField - sql %s args %s \n", queryInventory, argsInventory)

	if err != nil {
		return fmt.Errorf("updateInventory - u.Builder.Insert: %w", err)
	}

	if _, err = tx.Exec(ctx, queryInventory, argsInventory...); err != nil {
		return fmt.Errorf("updateInventory - tx.Exec: %w", err)
	}

	return nil
}

//func (u *InventoryRepository) GetByCustomerId(ctx context.Context, id string) (entity.Inventory, error) {
//	return u.getByField(ctx, "customer_id", id)
//}

//func (u *InventoryRepository) getByField(ctx context.Context, field, value string) (entity.Inventory, error) {
//	var err error
//	sql, args, _ := u.Builder.
//		Select("*").
//		From(inventoryTable).
//		Where(fmt.Sprintf("%v = ?", field), value).
//		ToSql()
//	log.Printf("InventoryRepository - getByField - sql %s args %s \n", sql, args)
//
//	var output entity.Inventory
//	err = u.Cluster.QueryRow(ctx, sql, args...).Scan(
//		&output.Id,
//		&output.CustomerId,
//		&output.Type,
//		&output.Quantity,
//	)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return entity.Inventory{}, repoerrs.ErrNotFound
//		}
//		return entity.Inventory{}, fmt.Errorf(
//			"InventoryRepository - getByField %s - r.Cluster.QueryRow: %v",
//			field,
//			err,
//		)
//	}
//	return output, nil
//}

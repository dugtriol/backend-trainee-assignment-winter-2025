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

func (u *InventoryRepository) Add(ctx context.Context, inventory entity.Inventory) (entity.Inventory, error) {
	var err error
	tx, err := u.Cluster.Begin(ctx)
	if err != nil {
		return entity.Inventory{}, err
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

	var output entity.Inventory
	// Вычитаем баланс пользователя
	if err = u.buyMerch(ctx, tx, inventory.CustomerID, inventory.Type); err != nil {
		return entity.Inventory{}, err
	}

	// Обновляем инвентарь
	if output, err = u.updateInventory(ctx, tx, inventory.CustomerID, inventory.Type); err != nil {
		return entity.Inventory{}, err
	}

	return output, nil
}

func (u *InventoryRepository) buyMerch(
	ctx context.Context, tx pgx.Tx, customerID string, merchType string,
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
				squirrel.Eq{"id": customerID},
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
	ctx context.Context, tx pgx.Tx, customerID string, merchType string,
) (entity.Inventory, error) {
	queryInventory, argsInventory, err := u.Builder.Insert(inventoryTable).
		Columns("customer_id", "type", "quantity").
		Values(customerID, merchType, 1).
		Suffix("ON CONFLICT (customer_id, type) DO UPDATE SET quantity = inventories.quantity + 1").
		Suffix("RETURNING id, customer_id, type, quantity").
		ToSql()
	log.Printf("InventoryRepository - getByField - sql %s args %s \n", queryInventory, argsInventory)

	if err != nil {
		return entity.Inventory{}, fmt.Errorf("updateInventory - u.Builder.Insert: %w", err)
	}

	var output entity.Inventory
	err = tx.QueryRow(ctx, queryInventory, argsInventory...).Scan(
		&output.ID,
		&output.CustomerID,
		&output.Type,
		&output.Quantity,
	)

	if err != nil {
		return entity.Inventory{}, fmt.Errorf("updateInventory - tx.Exec: %w", err)
	}

	return output, nil
}

func (r *InventoryRepository) GetByUserID(ctx context.Context, userID string) ([]entity.Inventory, error) {
	query, args, err := r.Builder.
		Select("*").
		From(inventoryTable).
		Where(squirrel.Eq{"customer_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("InventoryRepo - GetByUserID - r.Builder: %v", err)
	}

	rows, err := r.Cluster.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("InventoryRepo - GetByUserID - r.Cluster.Query: %v", err)
	}
	defer rows.Close()

	var inventoryList []entity.Inventory
	for rows.Next() {
		var item entity.Inventory
		if err = rows.Scan(
			&item.ID,
			&item.CustomerID,
			&item.Type,
			&item.Quantity,
		); err != nil {
			return nil, fmt.Errorf("InventoryRepo - GetByUserID - rows.Scan: %v", err)
		}
		inventoryList = append(inventoryList, item)
	}

	return inventoryList, nil
}

package pgdb

import (
	"context"
	"errors"
	"fmt"

	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/internal/repo/repoerrs"
	"backend-trainee-assignment-winter-2025/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	transactionTable = "transactions"
)

type TransactionRepository struct {
	*postgres.Database
}

func NewTransactionRepository(db *postgres.Database) *TransactionRepository {
	return &TransactionRepository{db}
}

func (r *TransactionRepository) GetByUserID(ctx context.Context, userId string) ([]entity.Transaction, error) {
	query, args, err := r.Builder.
		Select("*").
		From(transactionTable).
		Where(
			squirrel.Or{
				squirrel.Eq{"from_user": userId},
				squirrel.Eq{"to_user": userId},
			},
		).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - GetByUserID - r.Builder: %v", err)
	}

	rows, err := r.Cluster.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo - GetByUserID - r.Cluster.Query: %v", err)
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var txn entity.Transaction
		if err = rows.Scan(
			&txn.Id,
			&txn.FromUser,
			&txn.ToUser,
			&txn.Amount,
		); err != nil {
			return nil, fmt.Errorf("TransactionRepo - GetByUserID - rows.Scan: %v", err)
		}
		transactions = append(transactions, txn)
	}

	return transactions, nil
}

func (u *TransactionRepository) Transfer(
	ctx context.Context, input entity.Transaction,
	isExist func(ctx context.Context, id string) (entity.User, error),
) error {
	var err error
	tx, err := u.Cluster.Begin(ctx)
	if err != nil {
		return fmt.Errorf("TransactionRepository - starting transaction: %w", err)
	}

	defer func() {
		var e error
		if err == nil {
			e = tx.Commit(ctx)
		} else {
			e = tx.Rollback(ctx)
		}

		if err == nil && e != nil {
			err = fmt.Errorf("TransactionRepository - finishing transaction: %w", e)
		}
	}()

	// Проверяем существования пользователя
	if _, err = isExist(ctx, input.ToUser); err != nil {
		return err
	}

	// Проверяем баланс отправителя и списываем деньги
	if err = u.withdrawAmount(ctx, tx, input.FromUser, input.Amount); err != nil {
		return err
	}

	// Начисляем деньги получателю
	if err = u.depositAmount(ctx, tx, input.ToUser, input.Amount); err != nil {
		return err
	}

	// Добавляем запись в таблицу транзакций
	if err = u.addTransaction(ctx, tx, input.FromUser, input.ToUser, input.Amount); err != nil {
		return err
	}

	return nil
}

func (u *TransactionRepository) withdrawAmount(ctx context.Context, tx pgx.Tx, fromUserId string, amount int) error {
	query, args, err := u.Builder.
		Update(userTable).
		Set("amount", squirrel.Expr("amount - ?", amount)).
		Where(
			squirrel.And{
				squirrel.Eq{"id": fromUserId},
				squirrel.GtOrEq{"amount": amount},
			},
		).
		Suffix("RETURNING amount").
		ToSql()
	if err != nil {
		return fmt.Errorf("withdrawAmount - building query: %w", err)
	}

	var updatedAmount int
	if err = tx.QueryRow(ctx, query, args...).Scan(&updatedAmount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerrs.ErrLowBalance
		}
		return fmt.Errorf("withdrawAmount - executing query: %w", err)
	}

	return nil
}

func (u *TransactionRepository) depositAmount(ctx context.Context, tx pgx.Tx, toUserId string, amount int) error {
	query, args, err := u.Builder.
		Update(userTable).
		Set("amount", squirrel.Expr("amount + ?", amount)).
		Where(squirrel.Eq{"id": toUserId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("depositAmount - building query: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("depositAmount - executing query: %w", err)
	}

	return nil
}

func (u *TransactionRepository) addTransaction(
	ctx context.Context, tx pgx.Tx, fromUserId, toUserId string, amount int,
) error {
	query, args, err := u.Builder.
		Insert(transactionTable).
		Columns("from_user", "to_user", "amount").
		Values(fromUserId, toUserId, amount).
		ToSql()
	if err != nil {
		return fmt.Errorf("addTransaction - building query: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("addTransaction - executing query: %w", err)
	}

	return nil
}

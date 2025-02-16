//go:build integration
// +build integration

package postgres

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"backend-trainee-assignment-winter-2025/config"
	"backend-trainee-assignment-winter-2025/pkg/postgres"
	"github.com/Masterminds/squirrel"
)

const (
	ConnTestDatabase = "postgres://postgres:password@localhost:5433/testshop"
)

type TestDatabase struct {
	DB *postgres.Database
	sync.Mutex
}

func NewTestDatabase(configPath string) *TestDatabase {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(fmt.Sprintf("cfg: %v", cfg))

	database, err := postgres.New(
		context.Background(),
		cfg.Conn,
		postgres.MaxPoolSize(cfg.MaxPoolSize),
		postgres.ConnAttempts(cfg.Database.ConnAttempts),
		postgres.ConnTimeout(cfg.Database.ConnTimeout),
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &TestDatabase{DB: database}
}

func (d *TestDatabase) SetUp(t *testing.T) {
	t.Helper()
	d.Lock()
	if err := d.Truncate(context.Background()); err != nil {
		panic(err)
	}

	if err := d.AddMerch(context.Background()); err != nil {
		panic(err)
	}
}

func (d *TestDatabase) TearDown(t *testing.T) {
	defer d.Unlock()
	// можно закомментировать вызов Truncate, чтобы
	// посмотреть результаты действий над БД
	//if err := d.Truncate(context.Background()); err != nil {
	//	panic(err)
	//}
}

func (d *TestDatabase) AddMerch(ctx context.Context) error {
	_, err := d.DB.Cluster.Exec(
		ctx, `
		INSERT INTO merch (name, price) VALUES 
		('t-shirt', 80),
		('cup', 20),
		('book', 50),
		('pen', 10),
		('powerbank', 200),
		('hoody', 300),
		('umbrella', 200),
		('socks', 10),
		('wallet', 50),
		('pink-hoody', 500)
	ON CONFLICT (name) DO NOTHING`,
	)
	if err != nil {
		return fmt.Errorf("failed to add merch items: %w", err)
	}
	return nil
}
func (d *TestDatabase) Truncate(ctx context.Context) error {
	// Запрос списка таблиц
	query, args, err := d.DB.Builder.Select("table_name").
		From("information_schema.tables").
		Where(squirrel.Eq{"table_schema": "public", "table_type": "BASE TABLE"}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := d.DB.Cluster.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, table)
	}

	if len(tables) == 0 {
		return fmt.Errorf("no tables found")
	}

	// Создание TRUNCATE запроса
	truncateQuery := fmt.Sprintf("Truncate table" + " " + strings.Join(tables, ", "))
	_, err = d.DB.Cluster.Exec(ctx, truncateQuery)
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	return nil
}

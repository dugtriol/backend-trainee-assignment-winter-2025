//go:build integration
// +build integration

package tests

import (
	"fmt"
	"log"

	"backend-trainee-assignment-winter-2025/tests/postgres"
	"github.com/joho/godotenv"
)

var (
	db *postgres.TestDatabase
)

const (
	configPath = "../config/config.yaml"
	envPath    = "../.env"
)

func init() {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %v", err))
	}

	db = postgres.NewTestDatabase(configPath)
}

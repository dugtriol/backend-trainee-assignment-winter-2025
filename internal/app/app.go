package app

import (
	"context"
	"fmt"

	"backend-trainee-assignment-winter-2025/config"
	"backend-trainee-assignment-winter-2025/internal/repo"
	"backend-trainee-assignment-winter-2025/pkg/postgres"
)

func Run(configPath string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(cfg)

	//logger
	log := setLogger(cfg.Level)
	log.Info("Init logger")

	//postgres
	database, err := postgres.New(ctx, cfg.Conn, postgres.MaxPoolSize(cfg.MaxPoolSize))
	if err != nil {
		fmt.Println(err.Error())
	}

	//repositories
	repos := repo.NewRepositories(database)
	dependencies := service.ServicesDependencies{Repos: repos}

}

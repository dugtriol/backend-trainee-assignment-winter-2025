package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"backend-trainee-assignment-winter-2025/config"
	v1 "backend-trainee-assignment-winter-2025/internal/handler/http/v1"
	"backend-trainee-assignment-winter-2025/internal/repo"
	"backend-trainee-assignment-winter-2025/internal/service"
	"backend-trainee-assignment-winter-2025/pkg/httpserver"
	"backend-trainee-assignment-winter-2025/pkg/postgres"
	"github.com/go-chi/chi/v5"
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
	database, err := postgres.New(
		ctx,
		cfg.Conn,
		postgres.MaxPoolSize(cfg.MaxPoolSize),
		postgres.ConnAttempts(cfg.Database.ConnAttempts),
		postgres.ConnTimeout(cfg.Database.ConnTimeout),
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	//repositories
	repos := repo.NewRepositories(database)
	dependencies := service.ServicesDependencies{Repos: repos}

	//services
	services := service.NewServices(dependencies)

	//handlers
	log.Info("Initializing handlers and routes...")
	router := chi.NewRouter()
	v1.NewRouter(ctx, log, router, services)

	// HTTP server
	log.Info("Starting http server...")
	log.Debug(fmt.Sprintf("Server port: %s", cfg.Port))
	httpServer := httpserver.New(
		router,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.Timeout),
		httpserver.WriteTimeout(cfg.HTTP.Timeout),
		httpserver.ShutdownTimeout(cfg.HTTP.ShutdownTimeout),
	)

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}

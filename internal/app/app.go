package app

import (
	"avito-test2024-spring/internal/config"
	"avito-test2024-spring/internal/controller"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/internal/server"
	"avito-test2024-spring/internal/service"
	"avito-test2024-spring/pkg/auth"
	"avito-test2024-spring/pkg/database/postgresql"
	"avito-test2024-spring/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// TODO rewrite initial script for DB; add users table with ids and their tags

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	logs := logger.InitLogs(cfg.Logger)

	logs.Info().Msg("Starting app")
	logs.Info().Interface("config", cfg).Msg("")

	dbPool := postgresql.NewConnectionPool(cfg.PostgreSQL, logs)
	logs.Info().Msg("Initialized connection pool DB")

	repos := repository.NewRepositories(dbPool)
	logs.Info().Msg("Initialized repos")

	tokenManager, err := auth.NewManager(cfg.JWT.SigningKey)
	if err != nil {
		logs.Error().Err(err).Msg("error occured while init of token manager")
	}
	logs.Info().Msg("Initialized tokenManager")

	services := service.NewServices(repos, tokenManager)
	logs.Info().Msg("Initialized services")

	handlers := controller.NewHandler(services.Banners, services.Tags, services.Features, services.Users, logs, tokenManager)
	logs.Info().Msg("Initialized handlers")

	srv := server.NewServer(cfg.HTTP, handlers.Init("local", cfg.HTTP.Port))
	go func() {
		if err := srv.Run(); err != nil {
			logs.Error().Err(err).Msg("error occurred while running http server")
		}
	}()

	logs.Info().Msg("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	dbPool.Close()

	logs.Info().Msg("End of app")
}

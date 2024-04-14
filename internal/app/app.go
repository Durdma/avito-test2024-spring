package app

import (
	"avito-test2024-spring/internal/config"
	"avito-test2024-spring/internal/controller"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/internal/server"
	"avito-test2024-spring/internal/service"
	"avito-test2024-spring/pkg/auth"
	cache2 "avito-test2024-spring/pkg/cache"
	"avito-test2024-spring/pkg/database/postgresql"
	"avito-test2024-spring/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title Avito Banners API
// @version 1.0
// @description API для управления баннерами

// @host localhost:8080
// @BasePath /api/v1

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	logs := logger.NewLogs(cfg.Logger)

	logs.Logger.Info().Msg("Starting app")
	logs.Logger.Info().Interface("config", cfg).Msg("")

	cache := cache2.NewRedisCache(cfg.Cache)
	logs.Logger.Info().Msg("Initialized connection pool Cache")

	dbPool := postgresql.NewConnectionPool(cfg.PostgreSQL, logs)
	logs.Logger.Info().Msg("Initialized connection pool DB")

	repos := repository.NewRepositories(dbPool)
	logs.Logger.Info().Msg("Initialized repos")

	tokenManager, err := auth.NewManager(cfg.JWT.SigningKey)
	if err != nil {
		logs.Logger.Error().Err(err).Msg("error occured while init of token manager")
	}
	logs.Logger.Info().Msg("Initialized tokenManager")

	services := service.NewServices(repos, tokenManager, cache)
	logs.Logger.Info().Msg("Initialized services")

	handlers := controller.NewHandler(services.Banners, services.Tags, services.Features, services.Users, logs, tokenManager, cache)
	logs.Logger.Info().Msg("Initialized handlers")

	srv := server.NewServer(cfg.HTTP, handlers.Init("localhost", cfg.HTTP.Port))
	go func() {
		if err := srv.Run(); err != nil {
			logs.Logger.Error().Err(err).Msg("error occurred while running http server")
		}
	}()

	logs.Logger.Info().Msg("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	dbPool.Close()

	logs.Logger.Info().Msg("End of app")
}

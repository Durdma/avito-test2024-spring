package app

import (
	"avito-test2024-spring/internal/config"
	"avito-test2024-spring/pkg/logger"
	"log"
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
	logs.Info().Msg("End of app")
}

package app

import (
	"avito-test2024-spring/internal/config"
	"avito-test2024-spring/pkg/logger"
)

func Run(configPath string) {
	logs := logger.InitLogs("../../pkg/logger/logger.json", 5, 3, 30)

	logs.Info().Msg("Starting App")

	cfg, err := config.Init(configPath)
	if err != nil {
		logs.Error().Msg("Error while reading config")
		return
	}
	logs.Info().Interface("config", cfg).Msg("")
	logs.Info().Msg("End of app")
}

package postgresql

import (
	"avito-test2024-spring/internal/config"
	"avito-test2024-spring/pkg/database/postgresql/dbscripts"
	"avito-test2024-spring/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const timeout = 10 * time.Second

func NewConnectionPool(cfg config.PostgreSQLConfig, logs *logger.Logs) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, getConnectionString(cfg))
	if err != nil {
		logs.Logger.Error().Msg("error while connecting to DB")
		return nil
	}

	err = pool.Ping(context.Background())
	if err != nil {
		logs.Logger.Error().Msg("error while testing DB connection")
		logs.Logger.Error().Msg(err.Error())
		return nil
	}

	_, err = pool.Exec(context.Background(), dbscripts.Create)
	if err != nil {
		logs.Logger.Error().Msg("error while creating DB scheme")
		logs.Logger.Error().Msg(err.Error())
		return nil
	}

	return pool
}

func getConnectionString(cfg config.PostgreSQLConfig) string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

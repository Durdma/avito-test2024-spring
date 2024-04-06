package postgresql

import (
	"avito-test2024-spring/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"time"
)

const timeout = 10 * time.Second

func NewConnectionPool(cfg config.PostgreSQLConfig, logs zerolog.Logger) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, getConnectionString(cfg)) // TODO implement connection string
	if err != nil {
		logs.Error().Msg("error while connecting to DB")
		return nil
	}

	err = pool.Ping(context.Background())
	if err != nil {
		logs.Error().Msg("error while testing DB connection")
		return nil
	}

	return pool
}

// TODO заменить на переменные окружения
func getConnectionString(cfg config.PostgreSQLConfig) string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

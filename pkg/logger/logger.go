package logger

import (
	"avito-test2024-spring/internal/config"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func InitLogs(cfg config.LoggerConfig) zerolog.Logger {
	var writers []io.Writer

	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	writers = append(writers, &lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   true,
	})

	mw := io.MultiWriter(writers...)

	return zerolog.New(mw).With().Caller().Timestamp().Logger()
}

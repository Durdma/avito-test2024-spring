package logger

import (
	"avito-test2024-spring/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// TODO reformat logger and create methods for logging
// TODO use gin logger

type Logs struct {
	Logger zerolog.Logger
}

func NewLogs(cfg config.LoggerConfig) *Logs {
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

	return &Logs{
		Logger: zerolog.New(mw).With().Caller().Timestamp().Logger(),
	}
}

// Functions for logging api
func (l *Logs) Error(ctx *gin.Context, status int, err string) {
	l.Logger.Error().
		Str("method", ctx.Request.Method).
		Str("url", ctx.Request.RequestURI).
		Int("status_code", status).
		Msg(err)
}

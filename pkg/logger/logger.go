package logger

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// TODO add to config this params
func InitLogs(filename string, maxSize int, maxBackups int, maxAge int) zerolog.Logger {
	var writers []io.Writer

	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	writers = append(writers, &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	})

	mw := io.MultiWriter(writers...)

	return zerolog.New(mw).With().Caller().Timestamp().Logger()
}

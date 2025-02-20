package logs

import (
	"io"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func New(w io.Writer) *Logger {
	logger := slog.New(slog.NewJSONHandler(w, nil))
	return &Logger{logger: logger}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Error(msg string, err error, args ...interface{}) {
	args = append(args, "error", err.Error())
	l.logger.Error(msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{logger: l.logger.With(args...)}
}

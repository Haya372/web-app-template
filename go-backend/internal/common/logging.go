package common

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

type loggerImpl struct {
	logger slog.Logger
}

func (l *loggerImpl) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *loggerImpl) Info(ctx context.Context, msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *loggerImpl) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *loggerImpl) Error(ctx context.Context, msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func NewLogger() Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &loggerImpl{
		logger: *logger,
	}
}

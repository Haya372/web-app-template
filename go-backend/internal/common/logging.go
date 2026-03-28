package common

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
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
	l.logger.DebugContext(ctx, msg, withUserID(ctx, args)...)
}

func (l *loggerImpl) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, withUserID(ctx, args)...)
}

func (l *loggerImpl) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, withUserID(ctx, args)...)
}

func (l *loggerImpl) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, withUserID(ctx, args)...)
}

// withUserID prepends "userID" to args when a user ID is present in the context.
func withUserID(ctx context.Context, args []any) []any {
	userID := UserIDFromContext(ctx)
	if userID == "" {
		return args
	}

	return append([]any{"userID", userID}, args...)
}

func NewLogger() Logger {
	logger := otelslog.NewLogger("go-backend")

	return &loggerImpl{
		logger: *logger,
	}
}

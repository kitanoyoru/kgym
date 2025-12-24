package logger

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type contextKey string

const contextKeyLogger contextKey = "logger"

var ErrLoggerNotFound = errors.New("logger not found")

func FromContext(ctx context.Context) (*zap.Logger, error) {
	logger, ok := ctx.Value(contextKeyLogger).(*zap.Logger)
	if !ok {
		return zap.NewNop(), ErrLoggerNotFound
	}

	return logger, nil
}

func Inject(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

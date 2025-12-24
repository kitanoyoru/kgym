package logger

import (
	"context"

	"github.com/pkg/errors"
)

type contextKey string

const contextKeyLogger contextKey = "logger"

var ErrLoggerNotFound = errors.New("logger not found")

func FromContext(ctx context.Context) (*Logger, error) {
	logger, ok := ctx.Value(contextKeyLogger).(*Logger)
	if !ok {
		return nil, ErrLoggerNotFound
	}

	return logger, nil
}

func Inject(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

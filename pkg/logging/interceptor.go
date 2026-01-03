package logging

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog/log"
)

func NewInterceptorLogger(namespace string, service string) logging.Logger {
	l := log.With().Str("namespace", namespace).Str("service", service).Logger()

	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		logger := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			logger.Debug().Msg(msg)
		case logging.LevelInfo:
			logger.Info().Msg(msg)
		case logging.LevelWarn:
			logger.Warn().Msg(msg)
		case logging.LevelError:
			logger.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

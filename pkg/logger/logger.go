package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New(opts ...option) (*Logger, error) {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	var (
		zapLogger *zap.Logger
		err       error
	)

	if cfg.dev {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}

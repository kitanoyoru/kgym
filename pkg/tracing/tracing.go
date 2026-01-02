package tracing

import (
	"context"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/kitanoyoru/kgym/internal/apps/file/pkg/validator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type FlushFunc func(ctx context.Context) error

const (
	DefaultFlushTimeout = time.Second * 5
)

type Config struct {
	DSN   string `env:"SENTRY_DSN" validate:"required"`
	Debug bool   `env:"SENTRY_DEBUG" envDefault:"true"`

	ServiceName string `env:"SERVICE_NAME" validate:"required"`
}

func ConfigFromEnv(ctx context.Context) (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	if err := validator.Validate.StructCtx(ctx, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func Init(ctx context.Context, cfg Config) (FlushFunc, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.DSN,
		Debug:            cfg.Debug,
		AttachStacktrace: true,
		EnableLogs:       true,
		SendDefaultPII:   true,
		SampleRate:       1.0,
	})
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	flushFunc := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, DefaultFlushTimeout)
		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			return err
		}
		sentry.FlushWithContext(ctx)

		return nil
	}

	return flushFunc, nil
}

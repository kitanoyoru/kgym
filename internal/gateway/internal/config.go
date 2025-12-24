package gateway

import (
	"context"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	HTTPPort       string `env:"KGYM_GATEWAY_HTTP_PORT" validate:"required"`
	GRPCEndpoint   string `env:"KGYM_GATEWAY_GRPC_ENDPOINT" validate:"required"`
	EnableCors     bool   `env:"KGYM_GATEWAY_ENABLE_CORS" validate:"required"`
	BodyLimit      int    `env:"KGYM_GATEWAY_BODY_LIMIT" validate:"required"`
	MaxGRPCMsgSize int    `env:"KGYM_GATEWAY_MAX_GRPC_MSG_SIZE" validate:"required"`
}

func ParseAndValidate(ctx context.Context, cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	validate := validator.New(
		validator.WithPrivateFieldValidation(),
		validator.WithRequiredStructEnabled(),
	)

	return validate.StructCtx(ctx, cfg)
}

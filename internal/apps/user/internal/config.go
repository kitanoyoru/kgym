package internal

import (
	"context"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/user/pkg/validator"
)

type Config struct {
	Grpc
	Cache
	Database
}

func (c Config) Validate(ctx context.Context) error {
	return pkgValidator.Validate.StructCtx(ctx, c)
}

type Grpc struct {
}

type Cache struct {
}

type Database struct {
	ConnectionString string `env:"KGYM_USER_DATABASE_CONNECTION_STRING" validate:"required"`
}

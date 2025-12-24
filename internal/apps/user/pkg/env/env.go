package env

import (
	"context"

	"github.com/caarlos0/env/v11"
	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/validator"
)

func ParseAndValidate[T any](ctx context.Context, s *T) error {
	if err := env.Parse(s); err != nil {
		return err
	}

	return validator.Validate.StructCtx(ctx, s)
}

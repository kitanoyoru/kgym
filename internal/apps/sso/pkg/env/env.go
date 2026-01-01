package env

import (
	"context"

	"github.com/caarlos0/env/v10"
)

func ParseAndValidate(ctx context.Context, cfg interface{}) error {
	return env.Parse(cfg)
}

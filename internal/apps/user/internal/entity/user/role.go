package user

import (
	"context"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/user/pkg/validator"
)

const (
	Admin   Role = "admin"
	Default Role = "default"
)

type Role string

func (r Role) Validate(ctx context.Context) error {
	return pkgValidator.Validate.VarCtx(ctx, r, "oneof=admin default")
}

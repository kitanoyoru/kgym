package user

import (
	"context"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/user/pkg/validator"
)

type User struct {
	ID       string `validate:"required,uuid"`
	Email    string `validate:"required,email"`
	Role     Role
	Username string `validate:"required,min=3,max=32"`
	Password string `validate:"required,min=8,max=64"`
}

func (u User) Validate(ctx context.Context) error {
	if err := u.Role.Validate(ctx); err != nil {
		return err
	}

	return pkgValidator.Validate.StructCtx(ctx, u)
}

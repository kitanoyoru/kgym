package user

import (
	"context"
	"time"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/user/pkg/validator"
)

type User struct {
	ID        string `validate:"required,uuid"`
	Email     string `validate:"required,email"`
	Role      Role
	Username  string `validate:"required,min=3,max=32"`
	Password  string `validate:"required,min=8,max=64"`
	AvatarURL string `validate:"required,url"`
	Mobile    string `validate:"required,e164"`
	FirstName string `validate:"required,min=1,max=32"`
	LastName  string `validate:"required,min=1,max=32"`
	BirthDate time.Time
}

func (u User) Validate(ctx context.Context) error {
	if err := u.Role.Validate(ctx); err != nil {
		return err
	}

	return pkgValidator.Validate.StructCtx(ctx, u)
}

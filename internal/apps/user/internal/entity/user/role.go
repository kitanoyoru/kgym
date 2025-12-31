package user

import (
	"context"
	"strings"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/user/pkg/validator"
	"github.com/pkg/errors"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func RoleFromString(role string) (Role, error) {
	switch strings.ToLower(role) {
	case "admin":
		return RoleAdmin, nil
	case "user":
		return RoleUser, nil
	default:
		return "", errors.New("invalid role")
	}
}

func (r Role) Validate(ctx context.Context) error {
	return pkgValidator.Validate.VarCtx(ctx, r, "oneof=admin user")
}

package token

import (
	"context"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/sso/pkg/validator"
	"github.com/pkg/errors"
)

type Type string

const (
	TypeRefresh Type = "refresh"
)

func TypeFromString(s string) (Type, error) {
	switch s {
	case "refresh":
		return TypeRefresh, nil
	default:
		return "", errors.New("invalid token type")
	}
}

func (t Type) Validate(ctx context.Context) error {
	return pkgValidator.Validate.VarCtx(ctx, t, "oneof=refresh")
}

func (t Type) String() string {
	return string(t)
}

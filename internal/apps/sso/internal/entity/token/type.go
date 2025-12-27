package token

import (
	"context"

	"github.com/pkg/errors"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/sso/pkg/validator"
)

type TokenType string

const (
	TokenTypeRefresh TokenType = "refresh"
)

func TokenTypeFromString(s string) (TokenType, error) {
	switch s {
	case "refresh":
		return TokenTypeRefresh, nil
	default:
		return "", errors.New("invalid token type")
	}
}

func (t TokenType) Validate(ctx context.Context) error {
	return pkgValidator.Validate.VarCtx(ctx, t, "oneof=refresh")
}

func (t TokenType) String() string {
	return string(t)
}

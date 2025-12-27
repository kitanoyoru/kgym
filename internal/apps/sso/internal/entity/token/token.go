package token

import (
	"context"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/sso/pkg/validator"
)

type Token struct {
	ID        string
	UserID    string
	TokenType TokenType
	Token     string
}

func (t Token) Validate(ctx context.Context) error {
	if err := t.TokenType.Validate(ctx); err != nil {
		return err
	}

	return pkgValidator.Validate.StructCtx(ctx, t)
}

package token

import (
	"context"
	"time"

	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/sso/pkg/validator"
)

type Token struct {
	ID        string
	Subject   string
	ClientID  string
	TokenType Type
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
}

func (t Token) Validate(ctx context.Context) error {
	if err := t.TokenType.Validate(ctx); err != nil {
		return err
	}

	return pkgValidator.Validate.StructCtx(ctx, t)
}

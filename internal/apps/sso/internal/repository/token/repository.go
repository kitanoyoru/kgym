package token

import (
	"context"

	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token/models/token"
)

type IRepository interface {
	Create(ctx context.Context, token tokenentity.Token) error
	GetByTokenHash(ctx context.Context, tokenHash string) (tokenmodel.Token, error)
	Revoke(ctx context.Context, tokenHash string) error
}

package postgres

import (
	"context"

	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/models/token"
)

type IRepository interface {
	Create(ctx context.Context, token tokenmodel.Token) error
	List(ctx context.Context, filters ...Filter) ([]tokenmodel.Token, error)
	Update(ctx context.Context, id, token string) error
	Delete(ctx context.Context, filters ...Filter) error
}

package postgres

import (
	"context"

	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/models/token"
)

type IRepository interface {
	Create(ctx context.Context, token tokenentity.Token) error
	List(ctx context.Context, filters ...Filter) ([]tokenmodel.Token, error)
	Update(ctx context.Context, id, token string) error
	Delete(ctx context.Context, filters ...Filter) error
}

package postgres

import (
	"context"

	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/models/user"
)

type IRepository interface {
	Create(ctx context.Context, user usermodel.User) error
	List(ctx context.Context, filters ...Filter) ([]usermodel.User, error)
	Delete(ctx context.Context, filters ...Filter) error
}

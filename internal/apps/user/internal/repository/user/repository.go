package repository

import (
	"context"

	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/user/models/user"
)

type IRepository interface {
	GetByID(ctx context.Context, id string) (usermodel.User, error)
	GetByEmail(ctx context.Context, email string) (usermodel.User, error)
	Create(ctx context.Context, user userentity.User) error
	DeleteByID(ctx context.Context, id string) error
}

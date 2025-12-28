package user

import (
	"context"

	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user/models"
)

type IRepository interface {
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

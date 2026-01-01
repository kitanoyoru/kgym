package user

import (
	"context"
	"time"

	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
)

type IService interface {
	Create(ctx context.Context, req CreateRequest) (CreateResponse, error)
	GetByID(ctx context.Context, id string) (userentity.User, error)
	GetByEmail(ctx context.Context, email string) (userentity.User, error)
	DeleteByID(ctx context.Context, id string) error
}

type (
	CreateRequest struct {
		Email     string
		Role      userentity.Role
		Username  string
		Password  string
		AvatarURL string
		Mobile    string
		FirstName string
		LastName  string
		BirthDate time.Time
	}

	CreateResponse struct {
		ID string
	}
)

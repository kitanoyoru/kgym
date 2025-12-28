package grpc

import (
	"context"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	userrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user"
	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user/models"
)

var _ userrepo.IRepository = (*Repository)(nil)

type Repository struct {
	client pb.UserServiceClient
}

func New(client pb.UserServiceClient) *Repository {
	return &Repository{
		client,
	}
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	request := &pb.GetUser_Request{
		Email: &email,
	}

	resp, err := r.client.GetUser(ctx, request)
	if err != nil {
		return models.User{}, err
	}

	role, err := models.RoleFromProto(resp.User.Role)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:        resp.User.Id,
		Email:     resp.User.Email,
		Role:      role,
		Username:  resp.User.Username,
		Password:  resp.User.Password,
		AvatarURL: resp.User.AvatarUrl,
		Mobile:    resp.User.Mobile,
		FirstName: resp.User.FirstName,
		LastName:  resp.User.LastName,
		BirthDate: resp.User.BirthDate.AsTime(),
		CreatedAt: resp.User.CreatedAt.AsTime(),
	}, nil
}

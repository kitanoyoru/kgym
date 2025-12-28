package user

import (
	"context"

	"github.com/google/uuid"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	userrepo "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository"
)

type Service struct {
	repo userrepo.IRepository
}

func New(repo userrepo.IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	user := userentity.User{
		Email:     req.Email,
		Role:      req.Role,
		Username:  req.Username,
		Password:  req.Password,
		AvatarURL: req.AvatarURL,
		Mobile:    req.Mobile,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
	}

	user.ID = uuid.NewString()

	if err := s.repo.Create(ctx, user); err != nil {
		return CreateResponse{}, err
	}

	return CreateResponse{
		ID: user.ID,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (userentity.User, error) {
	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return userentity.User{}, err
	}

	role, err := model.Role.ToEntity()
	if err != nil {
		return userentity.User{}, err
	}

	return userentity.User{
		ID:        model.ID,
		Email:     model.Email,
		Role:      role,
		Username:  model.Username,
		Password:  model.Password,
		AvatarURL: model.AvatarURL,
		Mobile:    model.Mobile,
		FirstName: model.FirstName,
		LastName:  model.LastName,
	}, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (userentity.User, error) {
	model, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return userentity.User{}, err
	}

	role, err := model.Role.ToEntity()
	if err != nil {
		return userentity.User{}, err
	}

	return userentity.User{
		ID:        model.ID,
		Email:     model.Email,
		Role:      role,
		Username:  model.Username,
		Password:  model.Password,
		AvatarURL: model.AvatarURL,
		Mobile:    model.Mobile,
		FirstName: model.FirstName,
		LastName:  model.LastName,
	}, nil
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	return s.repo.DeleteByID(ctx, id)
}

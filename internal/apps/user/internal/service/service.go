package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/models/user"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/postgres"
	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics"
	servicemetrics "github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics/service"
)

type Service struct {
	repository postgres.IRepository
}

func New(repository postgres.IRepository) *Service {
	return &Service{
		repository,
	}
}

type CreateUserRequest struct {
	Email     string
	Role      string
	Username  string
	Password  string
	AvatarURL string
	Mobile    string
	FirstName string
	LastName  string
	BirthDate time.Time
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (string, error) {
	metrics.GlobalRegistry.GetMetric(servicemetrics.UserCreatedMetricName).Counter.WithLabelValues().Inc()

	userEntity := userentity.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Role:      userentity.Role(req.Role),
		Username:  req.Username,
		Password:  req.Password,
		AvatarURL: req.AvatarURL,
		Mobile:    req.Mobile,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
	}

	if err := userEntity.Validate(ctx); err != nil {
		return "", err
	}

	userModel := usermodel.User{
		ID:        userEntity.ID,
		Email:     userEntity.Email,
		Role:      usermodel.Role(userEntity.Role),
		Username:  userEntity.Username,
		Password:  userEntity.Password,
		AvatarURL: userEntity.AvatarURL,
		Mobile:    userEntity.Mobile,
		FirstName: userEntity.FirstName,
		LastName:  userEntity.LastName,
		BirthDate: userEntity.BirthDate,
	}

	err := s.repository.Create(ctx, userModel)
	if err != nil {
		return "", err
	}

	return userEntity.ID, nil
}

func (s *Service) List(ctx context.Context, options ...Option) ([]userentity.User, error) {
	metrics.GlobalRegistry.GetMetric(servicemetrics.UserListMetricName).Counter.WithLabelValues().Inc()

	var opts Options
	for _, option := range options {
		option(&opts)
	}

	var filters []postgres.Filter
	if opts.Email != nil {
		filters = append(filters, postgres.WithEmail(*opts.Email))
	}
	if opts.Username != nil {
		filters = append(filters, postgres.WithUsername(*opts.Username))
	}
	if opts.Role != nil {
		filters = append(filters, postgres.WithRole(usermodel.Role(*opts.Role)))
	}

	userModels, err := s.repository.List(ctx, filters...)
	if err != nil {
		return nil, err
	}

	users := make([]userentity.User, len(userModels))
	for i, model := range userModels {
		// TODO: Create serializer.go for service and move this to there
		users[i] = userentity.User{
			ID:        model.ID,
			Email:     model.Email,
			Role:      userentity.Role(model.Role),
			Username:  model.Username,
			Password:  model.Password,
			AvatarURL: model.AvatarURL,
			Mobile:    model.Mobile,
			FirstName: model.FirstName,
			LastName:  model.LastName,
			BirthDate: model.BirthDate,
		}
	}

	return users, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	metrics.GlobalRegistry.GetMetric(servicemetrics.UserDeletedMetricName).Counter.WithLabelValues().Inc()

	return s.repository.Delete(ctx, postgres.WithID(id))
}

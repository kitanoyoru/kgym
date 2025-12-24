package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/models/user"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/postgres"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/postgres/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"
)

type ServiceTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	mockRepo *mocks.MockIRepository
	service  *Service
	ctx      context.Context
}

func (s *ServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockRepo = mocks.NewMockIRepository(s.ctrl)
	s.service = New(s.mockRepo)
	s.ctx = context.Background()
}

func (s *ServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ServiceTestSuite) TestCreate() {
	s.Run("should create a user successfully", func() {
		req := CreateUserRequest{
			Email:    "test@example.com",
			Role:     "default",
			Username: "testuser",
			Password: "password123",
		}

		expectedUserModel := usermodel.User{
			ID:       "",
			Email:    req.Email,
			Role:     usermodel.RoleDefault,
			Username: req.Username,
			Password: req.Password,
		}

		s.mockRepo.EXPECT().
			Create(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, user usermodel.User) error {
				assert.NotEmpty(s.T(), user.ID)
				assert.Equal(s.T(), expectedUserModel.Email, user.Email)
				assert.Equal(s.T(), expectedUserModel.Role, user.Role)
				assert.Equal(s.T(), expectedUserModel.Username, user.Username)
				assert.Equal(s.T(), expectedUserModel.Password, user.Password)
				return nil
			})

		userID, err := s.service.Create(s.ctx, req)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), userID)
	})

	s.Run("should not create a user because of empty email", func() {
		req := CreateUserRequest{
			Email:    "",
			Role:     "default",
			Username: "testuser",
			Password: "password123",
		}

		userID, err := s.service.Create(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), userID)
	})

	s.Run("should not create a user because of invalid email", func() {
		req := CreateUserRequest{
			Email:    "invalid-email",
			Role:     "default",
			Username: "testuser",
			Password: "password123",
		}

		userID, err := s.service.Create(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), userID)
	})

	s.Run("should not create a user because of empty password", func() {
		req := CreateUserRequest{
			Email:    "test@example.com",
			Role:     "default",
			Username: "testuser",
			Password: "",
		}

		userID, err := s.service.Create(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), userID)
	})

	s.Run("should not create a user because of password too short", func() {
		req := CreateUserRequest{
			Email:    "test@example.com",
			Role:     "default",
			Username: "testuser",
			Password: "short",
		}

		userID, err := s.service.Create(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), userID)
	})

	s.Run("should not create a user because of username too short", func() {
		req := CreateUserRequest{
			Email:    "test@example.com",
			Role:     "default",
			Username: "ab",
			Password: "password123",
		}

		userID, err := s.service.Create(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), userID)
	})

	s.Run("should not create a user because of invalid role", func() {
		req := CreateUserRequest{
			Email:    "test@example.com",
			Role:     "invalid-role",
			Username: "testuser",
			Password: "password123",
		}

		userID, err := s.service.Create(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), userID)
	})

	s.Run("should not create a user because of repository error", func() {
		req := CreateUserRequest{
			Email:    "test@example.com",
			Role:     "default",
			Username: "testuser",
			Password: "password123",
		}

		expectedError := errors.New("repository error")
		s.mockRepo.EXPECT().
			Create(s.ctx, gomock.Any()).
			Return(expectedError)

		userID, err := s.service.Create(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		assert.Empty(s.T(), userID)
	})
}

func (s *ServiceTestSuite) TestList() {
	s.Run("should return list of users successfully with no filters", func() {
		userModels := []usermodel.User{
			{
				ID:        uuid.New().String(),
				Email:     "user1@example.com",
				Role:      usermodel.RoleDefault,
				Username:  "user1",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				Email:     "user2@example.com",
				Role:      usermodel.RoleAdmin,
				Username:  "user2",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx).
			Return(userModels, nil)

		users, err := s.service.List(s.ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 2)

		var ids []string
		for _, u := range users {
			ids = append(ids, u.ID)
		}
		assert.Contains(s.T(), ids, userModels[0].ID)
		assert.Contains(s.T(), ids, userModels[1].ID)
	})

	s.Run("should return list of users successfully with email filter", func() {
		email := "filter@example.com"
		userModels := []usermodel.User{
			{
				ID:        uuid.New().String(),
				Email:     email,
				Role:      usermodel.RoleDefault,
				Username:  "filteruser",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) ([]usermodel.User, error) {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.Email)
				assert.Equal(s.T(), email, *dbFilters.Email)
				return userModels, nil
			})

		users, err := s.service.List(s.ctx, WithEmail(email))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 1)
		assert.Equal(s.T(), email, users[0].Email)
	})

	s.Run("should return list of users successfully with username filter", func() {
		username := "filteruser"
		userModels := []usermodel.User{
			{
				ID:        uuid.New().String(),
				Email:     "filter@example.com",
				Role:      usermodel.RoleDefault,
				Username:  username,
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) ([]usermodel.User, error) {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.Username)
				assert.Equal(s.T(), username, *dbFilters.Username)
				return userModels, nil
			})

		users, err := s.service.List(s.ctx, WithUsername(username))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 1)
		assert.Equal(s.T(), username, users[0].Username)
	})

	s.Run("should return list of users successfully with role filter", func() {
		role := userentity.Admin
		userModels := []usermodel.User{
			{
				ID:        uuid.New().String(),
				Email:     "admin1@example.com",
				Role:      usermodel.RoleAdmin,
				Username:  "admin1",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				Email:     "admin2@example.com",
				Role:      usermodel.RoleAdmin,
				Username:  "admin2",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) ([]usermodel.User, error) {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.Role)
				assert.Equal(s.T(), usermodel.Role(role), *dbFilters.Role)
				return userModels, nil
			})

		users, err := s.service.List(s.ctx, WithRole(role))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 2)
		for _, u := range users {
			assert.Equal(s.T(), role, u.Role)
		}
	})

	s.Run("should return list of users successfully with multiple filters", func() {
		email := "multi@example.com"
		username := "multiuser"
		role := userentity.Default
		userModels := []usermodel.User{
			{
				ID:        uuid.New().String(),
				Email:     email,
				Role:      usermodel.RoleDefault,
				Username:  username,
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		s.mockRepo.EXPECT().
			List(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) ([]usermodel.User, error) {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.Email)
				assert.NotNil(s.T(), dbFilters.Username)
				assert.NotNil(s.T(), dbFilters.Role)
				assert.Equal(s.T(), email, *dbFilters.Email)
				assert.Equal(s.T(), username, *dbFilters.Username)
				assert.Equal(s.T(), usermodel.Role(role), *dbFilters.Role)
				return userModels, nil
			})

		users, err := s.service.List(s.ctx, WithEmail(email), WithUsername(username), WithRole(role))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 1)
		assert.Equal(s.T(), email, users[0].Email)
		assert.Equal(s.T(), username, users[0].Username)
		assert.Equal(s.T(), role, users[0].Role)
	})

	s.Run("should return empty list when no users found", func() {
		s.mockRepo.EXPECT().
			List(s.ctx).
			Return([]usermodel.User{}, nil)

		users, err := s.service.List(s.ctx)
		require.NoError(s.T(), err)
		assert.Empty(s.T(), users)
	})

	s.Run("should return error when repository fails", func() {
		expectedError := errors.New("repository error")
		s.mockRepo.EXPECT().
			List(s.ctx).
			Return(nil, expectedError)

		users, err := s.service.List(s.ctx)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		assert.Nil(s.T(), users)
	})
}

func (s *ServiceTestSuite) TestDelete() {
	s.Run("should delete a user successfully", func() {
		userID := uuid.New().String()

		s.mockRepo.EXPECT().
			Delete(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) error {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.ID)
				assert.Equal(s.T(), userID, *dbFilters.ID)
				return nil
			})

		err := s.service.Delete(s.ctx, userID)
		require.NoError(s.T(), err)
	})

	s.Run("should return error when repository fails", func() {
		userID := uuid.New().String()
		expectedError := errors.New("repository error")

		s.mockRepo.EXPECT().
			Delete(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, filters ...postgres.Filter) error {
				var dbFilters postgres.Filters
				for _, f := range filters {
					f(&dbFilters)
				}
				assert.NotNil(s.T(), dbFilters.ID)
				assert.Equal(s.T(), userID, *dbFilters.ID)
				return expectedError
			})

		err := s.service.Delete(s.ctx, userID)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

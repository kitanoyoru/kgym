package user

import (
	"context"
	"errors"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/google/uuid"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/mocks"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/models/user"
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
		req := CreateRequest{
			Email:     "test@example.com",
			Role:      userentity.RoleUser,
			Username:  "testuser",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		expectedUser := userentity.User{
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

		s.mockRepo.EXPECT().
			Create(s.ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, user userentity.User) error {
				assert.NotEmpty(s.T(), user.ID)
				assert.Regexp(s.T(), `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, user.ID)
				assert.Equal(s.T(), expectedUser.Email, user.Email)
				assert.Equal(s.T(), expectedUser.Role, user.Role)
				assert.Equal(s.T(), expectedUser.Username, user.Username)
				assert.Equal(s.T(), expectedUser.Password, user.Password)
				assert.Equal(s.T(), expectedUser.AvatarURL, user.AvatarURL)
				assert.Equal(s.T(), expectedUser.Mobile, user.Mobile)
				assert.Equal(s.T(), expectedUser.FirstName, user.FirstName)
				assert.Equal(s.T(), expectedUser.LastName, user.LastName)
				return nil
			})

		resp, err := s.service.Create(s.ctx, req)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), resp.ID)
		assert.Regexp(s.T(), `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, resp.ID)
	})

	s.Run("should return error when repository create fails", func() {
		req := CreateRequest{
			Email:     "test@example.com",
			Role:      userentity.RoleUser,
			Username:  "testuser",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		expectedErr := errors.New("repository error")
		s.mockRepo.EXPECT().
			Create(s.ctx, gomock.Any()).
			Return(expectedErr)

		resp, err := s.service.Create(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedErr, err)
		assert.Empty(s.T(), resp.ID)
	})
}

func (s *ServiceTestSuite) TestGetByID() {
	s.Run("should get user by id successfully", func() {
		userID := uuid.New().String()
		model := usermodel.User{
			ID:        userID,
			Email:     "test@example.com",
			Role:      usermodel.RoleUser,
			Username:  "testuser",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
		}

		s.mockRepo.EXPECT().
			GetByID(s.ctx, userID).
			Return(model, nil)

		user, err := s.service.GetByID(s.ctx, userID)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), userID, user.ID)
		assert.Equal(s.T(), model.Email, user.Email)
		assert.Equal(s.T(), userentity.RoleUser, user.Role)
		assert.Equal(s.T(), model.Username, user.Username)
		assert.Equal(s.T(), model.Password, user.Password)
		assert.Equal(s.T(), model.AvatarURL, user.AvatarURL)
		assert.Equal(s.T(), model.Mobile, user.Mobile)
		assert.Equal(s.T(), model.FirstName, user.FirstName)
		assert.Equal(s.T(), model.LastName, user.LastName)
	})

	s.Run("should return error when repository get by id fails", func() {
		userID := uuid.New().String()
		expectedErr := errors.New("user not found")

		s.mockRepo.EXPECT().
			GetByID(s.ctx, userID).
			Return(usermodel.User{}, expectedErr)

		user, err := s.service.GetByID(s.ctx, userID)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedErr, err)
		assert.Empty(s.T(), user.ID)
	})

	s.Run("should return error when role conversion fails", func() {
		userID := uuid.New().String()
		model := usermodel.User{
			ID:       userID,
			Email:    "test@example.com",
			Role:     usermodel.Role("invalid"),
			Username: "testuser",
			Password: "password123",
		}

		s.mockRepo.EXPECT().
			GetByID(s.ctx, userID).
			Return(model, nil)

		user, err := s.service.GetByID(s.ctx, userID)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), user.ID)
	})
}

func (s *ServiceTestSuite) TestGetByEmail() {
	s.Run("should get user by email successfully", func() {
		email := "test@example.com"
		model := usermodel.User{
			ID:        uuid.New().String(),
			Email:     email,
			Role:      usermodel.RoleAdmin,
			Username:  "testuser",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "Jane",
			LastName:  "Smith",
		}

		s.mockRepo.EXPECT().
			GetByEmail(s.ctx, email).
			Return(model, nil)

		user, err := s.service.GetByEmail(s.ctx, email)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), model.ID, user.ID)
		assert.Equal(s.T(), email, user.Email)
		assert.Equal(s.T(), userentity.RoleAdmin, user.Role)
		assert.Equal(s.T(), model.Username, user.Username)
		assert.Equal(s.T(), model.Password, user.Password)
		assert.Equal(s.T(), model.AvatarURL, user.AvatarURL)
		assert.Equal(s.T(), model.Mobile, user.Mobile)
		assert.Equal(s.T(), model.FirstName, user.FirstName)
		assert.Equal(s.T(), model.LastName, user.LastName)
	})

	s.Run("should return error when repository get by email fails", func() {
		email := "test@example.com"
		expectedErr := errors.New("user not found")

		s.mockRepo.EXPECT().
			GetByEmail(s.ctx, email).
			Return(usermodel.User{}, expectedErr)

		user, err := s.service.GetByEmail(s.ctx, email)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedErr, err)
		assert.Empty(s.T(), user.ID)
	})

	s.Run("should return error when role conversion fails", func() {
		email := "test@example.com"
		model := usermodel.User{
			ID:       uuid.New().String(),
			Email:    email,
			Role:     usermodel.Role("invalid"),
			Username: "testuser",
			Password: "password123",
		}

		s.mockRepo.EXPECT().
			GetByEmail(s.ctx, email).
			Return(model, nil)

		user, err := s.service.GetByEmail(s.ctx, email)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), user.ID)
	})
}

func (s *ServiceTestSuite) TestDeleteByID() {
	s.Run("should delete user by id successfully", func() {
		userID := uuid.New().String()

		s.mockRepo.EXPECT().
			DeleteByID(s.ctx, userID).
			Return(nil)

		err := s.service.DeleteByID(s.ctx, userID)
		assert.NoError(s.T(), err)
	})

	s.Run("should return error when repository delete fails", func() {
		userID := uuid.New().String()
		expectedErr := errors.New("delete failed")

		s.mockRepo.EXPECT().
			DeleteByID(s.ctx, userID).
			Return(expectedErr)

		err := s.service.DeleteByID(s.ctx, userID)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

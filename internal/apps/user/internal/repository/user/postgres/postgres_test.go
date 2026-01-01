package postgres

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/user/models/user"
	"github.com/kitanoyoru/kgym/internal/apps/user/migrations"
	postgresdb "github.com/kitanoyoru/kgym/pkg/database/postgres"
	"github.com/kitanoyoru/kgym/pkg/testing/integration/cockroachdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type RepositoryTestSuite struct {
	suite.Suite

	db        *pgxpool.Pool
	container *cockroachdb.CockroachDBContainer
}

func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	container, err := cockroachdb.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup test container")

	s.container = container

	s.db, err = postgresdb.New(ctx, postgresdb.Config{
		URI: container.URI,
	})
	require.NoError(s.T(), err, "failed to create postgres client")

	err = migrations.Up(ctx, "pgx", container.URI)
	require.NoError(s.T(), err, "failed to run migrations")
}

func (s *RepositoryTestSuite) TearDownSuite() {
	if s.container != nil {
		_ = s.container.Terminate(s.T().Context())
	}
	if s.db != nil {
		s.db.Close()
	}
}

func (s *RepositoryTestSuite) SetupTest() {
	ctx := context.Background()
	_, err := s.db.Exec(ctx, "DELETE FROM users")
	require.NoError(s.T(), err, "failed to clean users table")
}

func (s *RepositoryTestSuite) TearDownTest() {
	ctx := context.Background()
	_, err := s.db.Exec(ctx, "DELETE FROM users")
	require.NoError(s.T(), err, "failed to clean users table")
}

func (s *RepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should create a user successfully", func() {
		user := userentity.User{
			ID:        uuid.New().String(),
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

		err := repository.Create(ctx, user)
		assert.NoError(s.T(), err)

		retrievedUser, err := repository.GetByID(ctx, user.ID)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), user.ID, retrievedUser.ID)
		assert.Equal(s.T(), user.Email, retrievedUser.Email)
		assert.Equal(s.T(), usermodel.RoleUser, retrievedUser.Role)
		assert.Equal(s.T(), user.Username, retrievedUser.Username)
		assert.Equal(s.T(), user.Password, retrievedUser.Password)
	})

	s.Run("should not create a user because of duplicate email", func() {
		email := "duplicate@example.com"
		user1 := userentity.User{
			ID:        uuid.New().String(),
			Email:     email,
			Role:      userentity.RoleUser,
			Username:  "user1",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar1.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err := repository.Create(ctx, user1)
		require.NoError(s.T(), err)

		user2 := userentity.User{
			ID:        uuid.New().String(),
			Email:     email,
			Role:      userentity.RoleAdmin,
			Username:  "user2",
			Password:  "password456",
			AvatarURL: "https://example.com/avatar2.jpg",
			Mobile:    "+0987654321",
			FirstName: "Jane",
			LastName:  "Smith",
			BirthDate: carbon.CreateFromDateTime(1991, 2, 2, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err = repository.Create(ctx, user2)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a user because of duplicate username", func() {
		username := "duplicateuser"
		user1 := userentity.User{
			ID:        uuid.New().String(),
			Email:     "user1@example.com",
			Role:      userentity.RoleUser,
			Username:  username,
			Password:  "password123",
			AvatarURL: "https://example.com/avatar1.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err := repository.Create(ctx, user1)
		require.NoError(s.T(), err)

		user2 := userentity.User{
			ID:        uuid.New().String(),
			Email:     "user2@example.com",
			Role:      userentity.RoleAdmin,
			Username:  username,
			Password:  "password456",
			AvatarURL: "https://example.com/avatar2.jpg",
			Mobile:    "+0987654321",
			FirstName: "Jane",
			LastName:  "Smith",
			BirthDate: carbon.CreateFromDateTime(1991, 2, 2, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err = repository.Create(ctx, user2)
		assert.Error(s.T(), err)
	})
}

func (s *RepositoryTestSuite) TestGetByID() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should get user by id successfully", func() {
		user := userentity.User{
			ID:        uuid.New().String(),
			Email:     "getbyid@example.com",
			Role:      userentity.RoleUser,
			Username:  "getbyiduser",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err := repository.Create(ctx, user)
		require.NoError(s.T(), err)

		retrievedUser, err := repository.GetByID(ctx, user.ID)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), user.ID, retrievedUser.ID)
		assert.Equal(s.T(), user.Email, retrievedUser.Email)
		assert.Equal(s.T(), usermodel.RoleUser, retrievedUser.Role)
		assert.Equal(s.T(), user.Username, retrievedUser.Username)
		assert.Equal(s.T(), user.Password, retrievedUser.Password)
	})

	s.Run("should return error when user not found", func() {
		nonExistentID := uuid.New().String()
		_, err := repository.GetByID(ctx, nonExistentID)
		assert.Error(s.T(), err)
	})

	s.Run("should not return deleted user", func() {
		user := userentity.User{
			ID:        uuid.New().String(),
			Email:     "tobedeleted@example.com",
			Role:      userentity.RoleUser,
			Username:  "tobedeleted",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err := repository.Create(ctx, user)
		require.NoError(s.T(), err)

		err = repository.DeleteByID(ctx, user.ID)
		require.NoError(s.T(), err)

		_, err = repository.GetByID(ctx, user.ID)
		assert.Error(s.T(), err)
	})
}

func (s *RepositoryTestSuite) TestGetByEmail() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should get user by email successfully", func() {
		user := userentity.User{
			ID:        uuid.New().String(),
			Email:     "getbyemail@example.com",
			Role:      userentity.RoleAdmin,
			Username:  "getbyemailuser",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "Jane",
			LastName:  "Smith",
			BirthDate: carbon.CreateFromDateTime(1991, 2, 2, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err := repository.Create(ctx, user)
		require.NoError(s.T(), err)

		retrievedUser, err := repository.GetByEmail(ctx, user.Email)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), user.ID, retrievedUser.ID)
		assert.Equal(s.T(), user.Email, retrievedUser.Email)
		assert.Equal(s.T(), usermodel.RoleAdmin, retrievedUser.Role)
		assert.Equal(s.T(), user.Username, retrievedUser.Username)
		assert.Equal(s.T(), user.Password, retrievedUser.Password)
	})

	s.Run("should return error when user not found by email", func() {
		nonExistentEmail := "nonexistent@example.com"
		_, err := repository.GetByEmail(ctx, nonExistentEmail)
		assert.Error(s.T(), err)
	})

	s.Run("should not return deleted user by email", func() {
		user := userentity.User{
			ID:        uuid.New().String(),
			Email:     "tobedeletedbyemail@example.com",
			Role:      userentity.RoleUser,
			Username:  "tobedeletedbyemail",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err := repository.Create(ctx, user)
		require.NoError(s.T(), err)

		err = repository.DeleteByID(ctx, user.ID)
		require.NoError(s.T(), err)

		_, err = repository.GetByEmail(ctx, user.Email)
		assert.Error(s.T(), err)
	})
}

func (s *RepositoryTestSuite) TestDeleteByID() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should delete a user successfully", func() {
		user := userentity.User{
			ID:        uuid.New().String(),
			Email:     "todelete@example.com",
			Role:      userentity.RoleUser,
			Username:  "todelete",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err := repository.Create(ctx, user)
		require.NoError(s.T(), err)

		retrievedUser, err := repository.GetByID(ctx, user.ID)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), user.ID, retrievedUser.ID)

		err = repository.DeleteByID(ctx, user.ID)
		assert.NoError(s.T(), err)

		_, err = repository.GetByID(ctx, user.ID)
		assert.Error(s.T(), err)
	})

	s.Run("should not error when deleting non-existent user", func() {
		nonExistentID := uuid.New().String()
		err := repository.DeleteByID(ctx, nonExistentID)
		assert.NoError(s.T(), err)
	})

	s.Run("should allow deleting already deleted user", func() {
		user := userentity.User{
			ID:        uuid.New().String(),
			Email:     "doubledelete@example.com",
			Role:      userentity.RoleUser,
			Username:  "doubledelete",
			Password:  "password123",
			AvatarURL: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime(),
		}

		err := repository.Create(ctx, user)
		require.NoError(s.T(), err)

		err = repository.DeleteByID(ctx, user.ID)
		require.NoError(s.T(), err)

		err = repository.DeleteByID(ctx, user.ID)
		assert.NoError(s.T(), err)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

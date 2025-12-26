package postgres

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/models/user"
	"github.com/kitanoyoru/kgym/internal/apps/user/migrations"
	"github.com/kitanoyoru/kgym/pkg/database/postgres"
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

	s.db, err = postgres.New(ctx, postgres.Config{
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
	err := migrations.Up(ctx, "pgx", s.container.URI)
	require.NoError(s.T(), err, "failed to run migrations up")
}

func (s *RepositoryTestSuite) TearDownTest() {
	ctx := context.Background()
	err := migrations.Down(ctx, "pgx", s.container.URI)
	require.NoError(s.T(), err, "failed to run migrations down")
}

func (s *RepositoryTestSuite) TestCreate() {
	repository := New(s.db)
	ctx := context.Background()

	s.Run("should create a user successfully", func() {
		user := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "test@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "testuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, user)
		assert.NoError(s.T(), err)

		users, err := repository.List(ctx, WithEmail(user.Email))
		require.NoError(s.T(), err)
		require.Len(s.T(), users, 1)
		assert.Equal(s.T(), user.Email, users[0].Email)
		assert.Equal(s.T(), user.Username, users[0].Username)
	})

	s.Run("should not create a user because of non-unique email", func() {
		user1 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "duplicate@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "user1",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, user1)
		require.NoError(s.T(), err)

		user2 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "duplicate@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "user2",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err = repository.Create(ctx, user2)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a user because of empty email", func() {
		user := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "",
			Role:      usermodel.RoleDefault,
			Username:  "testuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, user)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a user because of email length is greater than 255 characters", func() {
		longEmail := make([]byte, 256)
		for i := range longEmail {
			longEmail[i] = 'a'
		}

		user := usermodel.User{
			ID:        uuid.New().String(),
			Email:     string(longEmail) + "@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "testuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, user)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a user because of empty password", func() {
		user := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "test@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "testuser",
			Password:  "",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, user)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a user because of password length is greater than 255 characters", func() {
		longPassword := make([]byte, 256)
		for i := range longPassword {
			longPassword[i] = 'a'
		}

		user := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "test@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "testuser",
			Password:  string(longPassword),
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, user)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a user because of non-unique username", func() {
		user1 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "user1@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "duplicateuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, user1)
		require.NoError(s.T(), err)

		user2 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "user2@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "duplicateuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err = repository.Create(ctx, user2)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a user because of username length is greater than 255 characters", func() {
		longUsername := make([]byte, 256)
		for i := range longUsername {
			longUsername[i] = 'a'
		}

		user := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "test@example.com",
			Role:      usermodel.RoleDefault,
			Username:  string(longUsername),
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, user)
		assert.Error(s.T(), err)
	})
}

func (s *RepositoryTestSuite) TestList() {
	repository := New(s.db)
	ctx := context.Background()

	clearUsers := func() {
		_, err := s.db.Exec(ctx, "DELETE FROM users")
		require.NoError(s.T(), err)
	}

	s.Run("shoud return list users successfully", func() {
		clearUsers()
		u1 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "user1@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "user1",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		u2 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "user2@example.com",
			Role:      usermodel.RoleAdmin,
			Username:  "user2",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, u1)
		require.NoError(s.T(), err)
		err = repository.Create(ctx, u2)
		require.NoError(s.T(), err)

		users, err := repository.List(ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 2)

		var ids []string
		for _, u := range users {
			ids = append(ids, u.ID)
		}
		assert.Contains(s.T(), ids, u1.ID)
		assert.Contains(s.T(), ids, u2.ID)
	})

	s.Run("should return all list of users because of no filters", func() {
		clearUsers()
		u1 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "all1@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "alluser1",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		u2 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "all2@example.com",
			Role:      usermodel.RoleAdmin,
			Username:  "alluser2",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		_ = repository.Create(ctx, u1)
		_ = repository.Create(ctx, u2)

		users, err := repository.List(ctx)
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 2)

		var ids []string
		for _, u := range users {
			ids = append(ids, u.ID)
		}
		assert.Contains(s.T(), ids, u1.ID)
		assert.Contains(s.T(), ids, u2.ID)
	})

	s.Run("should return empty list of users because of no users", func() {
		clearUsers()
		users, err := repository.List(ctx)
		require.NoError(s.T(), err)
		assert.Empty(s.T(), users)
	})

	s.Run("should return list of users because of email filter", func() {
		clearUsers()
		u1 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "filtremail1@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "auser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		u2 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "filtremail2@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "buser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, u1)
		require.NoError(s.T(), err)
		err = repository.Create(ctx, u2)
		require.NoError(s.T(), err)

		users, err := repository.List(ctx, WithEmail(u1.Email))
		require.NoError(s.T(), err)

		assert.Len(s.T(), users, 1)
		assert.Equal(s.T(), u1.Email, users[0].Email)
		assert.Equal(s.T(), u1.Username, users[0].Username)
	})

	s.Run("should return list of users because of username filter", func() {
		clearUsers()
		u1 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "username1@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "uniqueuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		u2 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "username2@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "otheruser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		_ = repository.Create(ctx, u1)
		_ = repository.Create(ctx, u2)

		users, err := repository.List(ctx, WithUsername(u1.Username))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 1)
		assert.Equal(s.T(), u1.Username, users[0].Username)
		assert.Equal(s.T(), u1.Email, users[0].Email)
	})

	s.Run("should return list of users because of role filter", func() {
		clearUsers()
		u1 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "roleuser1@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "roleuser1",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		u2 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "roleuser2@example.com",
			Role:      usermodel.RoleAdmin,
			Username:  "roleuser2",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		u3 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "roleuser3@example.com",
			Role:      usermodel.RoleAdmin,
			Username:  "roleuser3",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		_ = repository.Create(ctx, u1)
		_ = repository.Create(ctx, u2)
		_ = repository.Create(ctx, u3)

		users, err := repository.List(ctx, WithRole(usermodel.RoleAdmin))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 2)
		for _, u := range users {
			assert.Equal(s.T(), usermodel.RoleAdmin, u.Role)
		}
	})
}

func (s *RepositoryTestSuite) TestDelete() {
	repository := New(s.db)
	ctx := context.Background()

	clearUsers := func() {
		_, _ = s.db.Exec(ctx, "DELETE FROM users")
	}

	s.Run("should delete a user successfully", func() {
		clearUsers()

		u := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "deleteuser@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "deleteuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, u)
		require.NoError(s.T(), err)

		// Delete by ID
		err = repository.Delete(ctx, WithID(u.ID))
		require.NoError(s.T(), err)

		users, err := repository.List(ctx, WithID(u.ID))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 0)
	})

	s.Run("should delete a user because of id filter", func() {
		clearUsers()

		u := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "byid@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "byiduser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, u)
		require.NoError(s.T(), err)

		err = repository.Delete(ctx, WithID(u.ID))
		require.NoError(s.T(), err)

		users, err := repository.List(ctx, WithID(u.ID))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 0)
	})

	s.Run("should delete a user because of email filter", func() {
		clearUsers()

		u := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "byemail@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "byemailuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, u)
		require.NoError(s.T(), err)

		err = repository.Delete(ctx, WithEmail(u.Email))
		require.NoError(s.T(), err)

		users, err := repository.List(ctx, WithEmail(u.Email))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 0)
	})

	s.Run("should delete a user because of username filter", func() {
		clearUsers()

		u := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "byusername@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "byusernameuser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, u)
		require.NoError(s.T(), err)

		err = repository.Delete(ctx, WithUsername(u.Username))
		require.NoError(s.T(), err)

		users, err := repository.List(ctx, WithUsername(u.Username))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 0)
	})

	s.Run("should delete a user because of role filter", func() {
		clearUsers()

		u1 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "roledelete1@example.com",
			Role:      usermodel.RoleAdmin,
			Username:  "roledeleteuser1",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		u2 := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "roledelete2@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "roledeleteuser2",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}

		err := repository.Create(ctx, u1)
		require.NoError(s.T(), err)
		err = repository.Create(ctx, u2)
		require.NoError(s.T(), err)

		// Delete all users with RoleAdmin
		err = repository.Delete(ctx, WithRole(usermodel.RoleAdmin))
		require.NoError(s.T(), err)

		users, err := repository.List(ctx, WithRole(usermodel.RoleAdmin))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 0)

		usersDefault, err := repository.List(ctx, WithRole(usermodel.RoleDefault))
		require.NoError(s.T(), err)
		assert.Len(s.T(), usersDefault, 1)
		assert.Equal(s.T(), u2.Username, usersDefault[0].Username)
	})

	s.Run("should not delete a user because of no user", func() {
		clearUsers()

		nonexistentID := uuid.New().String()
		err := repository.Delete(ctx, WithID(nonexistentID))
		require.NoError(s.T(), err)
		// nothing breaks or panics. This is okay.
	})

	s.Run("should not delete a user because of no filters", func() {
		clearUsers()

		u := usermodel.User{
			ID:        uuid.New().String(),
			Email:     "nofilter@example.com",
			Role:      usermodel.RoleDefault,
			Username:  "nofilteruser",
			Password:  "password123",
			CreatedAt: carbon.Now().StdTime(),
			UpdatedAt: carbon.Now().StdTime(),
		}
		err := repository.Create(ctx, u)
		require.NoError(s.T(), err)

		// Should not delete any user if no filters are provided
		err = repository.Delete(ctx)
		require.NoError(s.T(), err)

		users, err := repository.List(ctx, WithID(u.ID))
		require.NoError(s.T(), err)
		assert.Len(s.T(), users, 1)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

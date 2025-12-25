package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	filemodel "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/models/file"
	"github.com/kitanoyoru/kgym/internal/apps/file/migrations"
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
	_, err := s.db.Exec(ctx, "DELETE FROM files")
	require.NoError(s.T(), err, "failed to clean files table")
}

func (s *RepositoryTestSuite) TearDownTest() {
	ctx := context.Background()
	_, err := s.db.Exec(ctx, "DELETE FROM files")
	require.NoError(s.T(), err, "failed to clean files table")
}

func (s *RepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should create a file successfully", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      "/test/path/image.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		assert.NoError(s.T(), err)

		files, err := repository.List(ctx, WithPath(file.Path))
		require.NoError(s.T(), err)
		require.Len(s.T(), files, 1)
		assert.Equal(s.T(), file.Path, files[0].Path)
		assert.Equal(s.T(), file.UserID, files[0].UserID)
		assert.Equal(s.T(), file.Size, files[0].Size)
		assert.Equal(s.T(), file.Extension, files[0].Extension)
		assert.Equal(s.T(), file.State, files[0].State)
	})

	s.Run("should not create a file because of non-unique path", func() {
		path := "/test/path/duplicate.png"
		file1 := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      path,
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file1)
		require.NoError(s.T(), err)

		file2 := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      path,
			Size:      2048,
			Extension: filemodel.ExtensionJPEG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = repository.Create(ctx, file2)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a file because of empty user id", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    "",
			Path:      "/test/path/image.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a file because of empty size", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      "/test/path/image.png",
			Size:      0,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a file because of empty extension", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      "/test/path/image.png",
			Size:      1024,
			Extension: "",
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a file because of invalid extension", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      "/test/path/image.xyz",
			Size:      1024,
			Extension: filemodel.Extension("xyz"),
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a file because of invalid size", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      "/test/path/image.png",
			Size:      -1,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a file because of invalid user id", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    "invalid-uuid",
			Path:      "/test/path/image.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		assert.Error(s.T(), err)
	})
}

func (s *RepositoryTestSuite) TestList() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should list files successfully", func() {
		userID := uuid.New().String()
		file1 := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    userID,
			Path:      "/test/path/image1.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		file2 := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    userID,
			Path:      "/test/path/image2.jpg",
			Size:      2048,
			Extension: filemodel.ExtensionJPEG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file1)
		require.NoError(s.T(), err)

		err = repository.Create(ctx, file2)
		require.NoError(s.T(), err)

		files, err := repository.List(ctx, WithUserID(userID))
		require.NoError(s.T(), err)
		require.Len(s.T(), files, 2)

		paths := make(map[string]bool)
		for _, f := range files {
			paths[f.Path] = true
		}
		assert.True(s.T(), paths[file1.Path])
		assert.True(s.T(), paths[file2.Path])
	})

	s.Run("should not list files because of invalid user id", func() {
		files, err := repository.List(ctx, WithUserID("invalid-uuid"))
		assert.Error(s.T(), err)
		assert.Nil(s.T(), files)
	})

	s.Run("should list files with empty path filter", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      "/test/path/valid.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		require.NoError(s.T(), err)

		files, err := repository.List(ctx, WithPath(""))
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), files)
	})

	s.Run("should not list files because of invalid size", func() {
		size := int64(-1)
		files, err := repository.List(ctx, WithSize(size))
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), files)
	})

	s.Run("should not list files because they are deleted", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      "/test/path/deleted.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		require.NoError(s.T(), err)

		err = repository.Delete(ctx, WithID(file.ID))
		require.NoError(s.T(), err)

		files, err := repository.List(ctx, WithID(file.ID))
		require.NoError(s.T(), err)
		assert.Empty(s.T(), files)
	})
}

func (s *RepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should delete a file successfully", func() {
		file := filemodel.File{
			ID:        uuid.New().String(),
			UserID:    uuid.New().String(),
			Path:      "/test/path/to-delete.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StatePending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.Create(ctx, file)
		require.NoError(s.T(), err)

		files, err := repository.List(ctx, WithID(file.ID))
		require.NoError(s.T(), err)
		require.Len(s.T(), files, 1)

		err = repository.Delete(ctx, WithID(file.ID))
		assert.NoError(s.T(), err)

		files, err = repository.List(ctx, WithID(file.ID))
		require.NoError(s.T(), err)
		assert.Empty(s.T(), files)
	})

	s.Run("should not delete a file because file not found", func() {
		nonExistentID := uuid.New().String()
		err := repository.Delete(ctx, WithID(nonExistentID))
		assert.NoError(s.T(), err)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

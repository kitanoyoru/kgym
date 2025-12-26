package service

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/minio"
	filemodel "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/models/file"
	"github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/postgres"
	"github.com/kitanoyoru/kgym/internal/apps/file/migrations"
	pkgpostgres "github.com/kitanoyoru/kgym/pkg/database/postgres"
	"github.com/kitanoyoru/kgym/pkg/testing/integration/cockroachdb"
	pkgminio "github.com/kitanoyoru/kgym/pkg/testing/integration/minio"
	minioClient "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type ServiceTestSuite struct {
	suite.Suite

	ctx context.Context

	dbContainer    *cockroachdb.CockroachDBContainer
	minioContainer *pkgminio.MinioContainer
}

func (s *ServiceTestSuite) SetupSuite() {
	ctx := context.Background()
	s.ctx = ctx

	dbContainer, err := cockroachdb.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup database container")
	s.dbContainer = dbContainer

	err = migrations.Up(ctx, "pgx", dbContainer.URI)
	require.NoError(s.T(), err, "failed to run migrations")

	minioContainer, err := pkgminio.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup minio container")
	s.minioContainer = minioContainer

	client, err := minioClient.New(minioContainer.Endpoint, &minioClient.Options{
		Creds: credentials.NewStaticV4(minioContainer.AccessKey, minioContainer.SecretKey, ""),
	})
	require.NoError(s.T(), err, "failed to create minio client")

	err = client.MakeBucket(ctx, "user-avatar", minioClient.MakeBucketOptions{})
	require.NoError(s.T(), err, "failed to create bucket")
}

func (s *ServiceTestSuite) TearDownSuite() {
	if s.dbContainer != nil {
		_ = s.dbContainer.Terminate(s.T().Context())
	}
	if s.minioContainer != nil {
		_ = s.minioContainer.Terminate(s.T().Context())
	}
}

func (s *ServiceTestSuite) SetupTest() {
	db, err := pkgpostgres.New(s.ctx, pkgpostgres.Config{
		URI: s.dbContainer.URI,
	})
	require.NoError(s.T(), err, "failed to create postgres client")

	_, err = db.Exec(s.ctx, "DELETE FROM files")
	require.NoError(s.T(), err, "failed to clean files table")

	db.Close()
}

func (s *ServiceTestSuite) createService() (*Service, *postgres.Repository, *minio.Repository, *pgxpool.Pool) {
	db, err := pkgpostgres.New(s.ctx, pkgpostgres.Config{
		URI: s.dbContainer.URI,
	})
	require.NoError(s.T(), err, "failed to create postgres client")

	client, err := minioClient.New(s.minioContainer.Endpoint, &minioClient.Options{
		Creds: credentials.NewStaticV4(s.minioContainer.AccessKey, s.minioContainer.SecretKey, ""),
	})
	require.NoError(s.T(), err, "failed to create minio client")

	postgresRepo := postgres.New(db)
	minioRepo, err := minio.New(client)
	require.NoError(s.T(), err, "failed to create minio repository")

	service := New(Config{
		Buckets: map[string]string{
			"user_avatar": "user-avatar",
		},
	}, minioRepo, postgresRepo)

	return service, postgresRepo, minioRepo, db
}

func (s *ServiceTestSuite) TestUpload() {
	s.Run("should upload a file successfully", func() {
		service, postgresRepo, _, db := s.createService()
		defer db.Close()

		userID := uuid.New().String()
		fileName := "avatar.png"
		contentType := "image/png"
		fileContent := []byte("test file content")
		reader := bytes.NewReader(fileContent)

		req := UploadRequest{
			UserID:      userID,
			Target:      "user_avatar",
			Name:        fileName,
			ContentType: contentType,
			Reader:      reader,
		}

		resp, err := service.Upload(s.ctx, req)
		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), resp.ID)

		file, err := postgresRepo.Get(s.ctx, resp.ID)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), userID, file.UserID)
		assert.Equal(s.T(), filemodel.ExtensionPNG, file.Extension)
		assert.Equal(s.T(), filemodel.StateCompleted, file.State)
		assert.Equal(s.T(), int64(len(fileContent)), file.Size)
	})

	s.Run("should return error when bucket not found", func() {
		_, _, minioRepo, db := s.createService()
		defer db.Close()

		serviceWithEmptyBucket := New(Config{
			Buckets: map[string]string{
				"user_avatar": "",
			},
		}, minioRepo, postgres.New(db))

		req := UploadRequest{
			UserID:      uuid.New().String(),
			Target:      "user_avatar",
			Name:        "avatar.png",
			ContentType: "image/png",
			Reader:      bytes.NewReader([]byte("test")),
		}

		resp, err := serviceWithEmptyBucket.Upload(s.ctx, req)
		assert.Error(s.T(), err)
		assert.Equal(s.T(), ErrBucketNotFound, err)
		assert.Empty(s.T(), resp.ID)
	})

	s.Run("should return error when file extension is invalid", func() {
		service, _, _, db := s.createService()
		defer db.Close()

		req := UploadRequest{
			UserID:      uuid.New().String(),
			Target:      "user_avatar",
			Name:        "avatar.invalid",
			ContentType: "image/png",
			Reader:      bytes.NewReader([]byte("test")),
		}

		resp, err := service.Upload(s.ctx, req)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "invalid extension"))
		assert.Empty(s.T(), resp.ID)
	})
}

func (s *ServiceTestSuite) TestGetURL() {
	s.Run("should get a file URL successfully", func() {
		service, _, _, db := s.createService()
		defer db.Close()

		userID := uuid.New().String()
		fileName := "avatar.png"
		contentType := "image/png"
		fileContent := []byte("test file content")
		reader := bytes.NewReader(fileContent)

		req := UploadRequest{
			UserID:      userID,
			Target:      "user_avatar",
			Name:        fileName,
			ContentType: contentType,
			Reader:      reader,
		}

		uploadResp, err := service.Upload(s.ctx, req)
		require.NoError(s.T(), err)

		url, err := service.GetURL(s.ctx, uploadResp.ID)
		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), url)
		assert.Contains(s.T(), url, "user-avatar")
	})

	s.Run("should return error when file not found in database", func() {
		service, _, _, db := s.createService()
		defer db.Close()

		fileID := uuid.New().String()

		url, err := service.GetURL(s.ctx, fileID)
		assert.Error(s.T(), err)
		assert.Empty(s.T(), url)
	})
}

func (s *ServiceTestSuite) TestDelete() {
	s.Run("should delete a file successfully", func() {
		service, postgresRepo, _, db := s.createService()
		defer db.Close()

		userID := uuid.New().String()
		fileName := "avatar.png"
		contentType := "image/png"
		fileContent := []byte("test file content")
		reader := bytes.NewReader(fileContent)

		req := UploadRequest{
			UserID:      userID,
			Target:      "user_avatar",
			Name:        fileName,
			ContentType: contentType,
			Reader:      reader,
		}

		uploadResp, err := service.Upload(s.ctx, req)
		require.NoError(s.T(), err)

		file, err := postgresRepo.Get(s.ctx, uploadResp.ID)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), file.Path)

		err = service.Delete(s.ctx, uploadResp.ID)
		assert.NoError(s.T(), err)

		_, err = postgresRepo.Get(s.ctx, uploadResp.ID)
		assert.Error(s.T(), err)
	})

	s.Run("should return error when file not found in database", func() {
		service, _, _, db := s.createService()
		defer db.Close()

		fileID := uuid.New().String()

		err := service.Delete(s.ctx, fileID)
		assert.Error(s.T(), err)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

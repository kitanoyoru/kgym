package minio

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	pkgminio "github.com/kitanoyoru/kgym/pkg/testing/integration/minio"
)

type RepositoryTestSuite struct {
	suite.Suite

	container *pkgminio.MinioContainer
}

func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	container, err := pkgminio.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup test container")

	s.container = container
}

func (s *RepositoryTestSuite) TearDownSuite() {
	if s.container != nil {
		_ = s.container.Terminate(s.T().Context())
	}
}

func (s *RepositoryTestSuite) TestUpload() {
	ctx := context.Background()

	minioClient, err := minio.New(s.container.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(s.container.AccessKey, s.container.SecretKey, ""),
	})
	require.NoError(s.T(), err, "failed to create minio client")

	repository := New(minioClient)

	s.Run("should upload a file successfully", func() {
		bucketName := "test-bucket-" + uuid.New().String()
		err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		require.NoError(s.T(), err, "failed to create bucket")

		fileName := "test-image.png"
		fileContent := "test file content for image"
		contentType := "image/png"

		req := UploadRequest{
			Bucket:      bucketName,
			Name:        fileName,
			ContentType: contentType,
			Reader:      strings.NewReader(fileContent),
		}

		resp, err := repository.Upload(ctx, req)
		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), resp.URL)
		assert.Equal(s.T(), "png", resp.Extension)
		assert.Equal(s.T(), int64(len(fileContent)), resp.Size)
		assert.Contains(s.T(), resp.URL, bucketName)
		assert.Contains(s.T(), resp.URL, fileName)
	})
}

func (s *RepositoryTestSuite) TestGetURL() {
	ctx := context.Background()

	minioClient, err := minio.New(s.container.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(s.container.AccessKey, s.container.SecretKey, ""),
	})
	require.NoError(s.T(), err, "failed to create minio client")

	repository := New(minioClient)

	s.Run("should get a file URL successfully", func() {
		bucketName := "test-bucket-" + uuid.New().String()
		err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		require.NoError(s.T(), err, "failed to create bucket")

		fileName := "test-document.pdf"
		fileContent := "test pdf content"
		contentType := "application/pdf"

		uploadReq := UploadRequest{
			Bucket:      bucketName,
			Name:        fileName,
			ContentType: contentType,
			Reader:      strings.NewReader(fileContent),
		}

		_, err = repository.Upload(ctx, uploadReq)
		require.NoError(s.T(), err, "failed to upload file")

		objects := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
		var objectName string
		for object := range objects {
			require.NoError(s.T(), object.Err, "failed to list objects")
			objectName = object.Key
			break
		}
		require.NotEmpty(s.T(), objectName, "object should exist")

		path := bucketName + "/" + objectName
		url, err := repository.GetURL(ctx, path)
		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), url)
		assert.Contains(s.T(), url, "response-content-disposition")
		assert.Contains(s.T(), url, objectName)
	})
}

func (s *RepositoryTestSuite) TestDelete() {
	ctx := context.Background()

	minioClient, err := minio.New(s.container.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(s.container.AccessKey, s.container.SecretKey, ""),
	})
	require.NoError(s.T(), err, "failed to create minio client")

	repository := New(minioClient)

	s.Run("should delete a file successfully", func() {
		bucketName := "test-bucket-" + uuid.New().String()
		err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		require.NoError(s.T(), err, "failed to create bucket")

		fileName := "test-file.txt"
		fileContent := "test content"
		contentType := "text/plain"

		uploadReq := UploadRequest{
			Bucket:      bucketName,
			Name:        fileName,
			ContentType: contentType,
			Reader:      strings.NewReader(fileContent),
		}

		_, err = repository.Upload(ctx, uploadReq)
		require.NoError(s.T(), err, "failed to upload file")

		objects := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
		var objectName string
		for object := range objects {
			require.NoError(s.T(), object.Err, "failed to list objects")
			objectName = object.Key
			break
		}
		require.NotEmpty(s.T(), objectName, "object should exist")

		_, err = minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
		require.NoError(s.T(), err, "file should exist before deletion")

		path := bucketName + "/" + objectName
		err = repository.Delete(ctx, path)
		assert.NoError(s.T(), err)

		_, err = minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
		assert.Error(s.T(), err, "file should not exist after deletion")
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

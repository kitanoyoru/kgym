package service

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/minio"
	miniomocks "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/minio/mocks"
	filemodel "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/models/file"
	postgresmocks "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/postgres/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"
)

type ServiceTestSuite struct {
	suite.Suite

	ctx  context.Context
	ctrl *gomock.Controller

	service *Service

	mockMinioRepo    *miniomocks.MockIRepository
	mockPostgresRepo *postgresmocks.MockIRepository
}

func (s *ServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockMinioRepo = miniomocks.NewMockIRepository(s.ctrl)
	s.mockPostgresRepo = postgresmocks.NewMockIRepository(s.ctrl)
	s.service = New(Config{
		Buckets: map[string]string{
			"user_avatar": "user-avatar",
		},
	}, s.mockMinioRepo, s.mockPostgresRepo)
	s.ctx = context.Background()
}

func (s *ServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ServiceTestSuite) TestUpload() {
	s.Run("should upload a file successfully", func() {
		userID := uuid.New().String()
		fileName := "avatar.png"
		contentType := "image/png"
		reader := bytes.NewReader([]byte("test file content"))

		s.mockPostgresRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, file filemodel.File) error {
				assert.Equal(s.T(), userID, file.UserID)
				assert.Equal(s.T(), "user-avatar/avatar.png", file.Path)
				assert.Equal(s.T(), filemodel.ExtensionPNG, file.Extension)
				assert.Equal(s.T(), filemodel.StatePending, file.State)
				return nil
			})

		s.mockMinioRepo.EXPECT().
			Upload(gomock.Any(), gomock.Any()).
			Return(minio.UploadResponse{}, nil)

		s.mockPostgresRepo.EXPECT().
			UpdateState(gomock.Any(), gomock.Any(), filemodel.StateCompleted).
			Return(nil)

		req := UploadRequest{
			UserID:      userID,
			Target:      "user_avatar",
			Name:        fileName,
			ContentType: contentType,
			Reader:      reader,
		}

		resp, err := s.service.Upload(s.ctx, req)
		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), resp.ID)
	})

	s.Run("should return error when bucket not found", func() {
		serviceWithEmptyBucket := New(Config{
			Buckets: map[string]string{
				"user_avatar": "",
			},
		}, s.mockMinioRepo, s.mockPostgresRepo)

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
		req := UploadRequest{
			UserID:      uuid.New().String(),
			Target:      "user_avatar",
			Name:        "avatar.invalid",
			ContentType: "image/png",
			Reader:      bytes.NewReader([]byte("test")),
		}

		resp, err := s.service.Upload(s.ctx, req)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "invalid file extension"))
		assert.Empty(s.T(), resp.ID)
	})

	s.Run("should return error when postgres create fails", func() {
		userID := uuid.New().String()
		expectedErr := errors.New("database error")

		s.mockPostgresRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		req := UploadRequest{
			UserID:      userID,
			Target:      "user_avatar",
			Name:        "avatar.png",
			ContentType: "image/png",
			Reader:      bytes.NewReader([]byte("test")),
		}

		resp, err := s.service.Upload(s.ctx, req)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "failed to create file record"))
		assert.Empty(s.T(), resp.ID)
	})

	s.Run("should return error when minio upload fails and update state to failed", func() {
		userID := uuid.New().String()
		var fileID string
		uploadErr := errors.New("minio upload error")

		s.mockPostgresRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, file filemodel.File) error {
				fileID = file.ID
				return nil
			})

		s.mockMinioRepo.EXPECT().
			Upload(gomock.Any(), gomock.Any()).
			Return(minio.UploadResponse{}, uploadErr)

		s.mockPostgresRepo.EXPECT().
			UpdateState(gomock.Any(), gomock.Any(), filemodel.StateFailed).
			DoAndReturn(func(ctx context.Context, id string, state filemodel.State) error {
				assert.Equal(s.T(), fileID, id)
				return nil
			})

		req := UploadRequest{
			UserID:      userID,
			Target:      "user_avatar",
			Name:        "avatar.png",
			ContentType: "image/png",
			Reader:      bytes.NewReader([]byte("test")),
		}

		resp, err := s.service.Upload(s.ctx, req)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "upload failed"))
		assert.Empty(s.T(), resp.ID)
	})

	s.Run("should return error when minio upload fails and state update also fails", func() {
		userID := uuid.New().String()
		var fileID string
		uploadErr := errors.New("minio upload error")
		updateErr := errors.New("state update error")

		s.mockPostgresRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, file filemodel.File) error {
				fileID = file.ID
				return nil
			})

		s.mockMinioRepo.EXPECT().
			Upload(gomock.Any(), gomock.Any()).
			Return(minio.UploadResponse{}, uploadErr)

		s.mockPostgresRepo.EXPECT().
			UpdateState(gomock.Any(), gomock.Any(), filemodel.StateFailed).
			DoAndReturn(func(ctx context.Context, id string, state filemodel.State) error {
				assert.Equal(s.T(), fileID, id)
				return updateErr
			})

		req := UploadRequest{
			UserID:      userID,
			Target:      "user_avatar",
			Name:        "avatar.png",
			ContentType: "image/png",
			Reader:      bytes.NewReader([]byte("test")),
		}

		resp, err := s.service.Upload(s.ctx, req)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "upload failed"))
		assert.True(s.T(), strings.Contains(err.Error(), "state update failed"))
		assert.Empty(s.T(), resp.ID)
	})

	s.Run("should return error when state update to completed fails", func() {
		userID := uuid.New().String()
		var fileID string
		updateErr := errors.New("state update error")

		s.mockPostgresRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, file filemodel.File) error {
				fileID = file.ID
				return nil
			})

		s.mockMinioRepo.EXPECT().
			Upload(gomock.Any(), gomock.Any()).
			Return(minio.UploadResponse{}, nil)

		s.mockPostgresRepo.EXPECT().
			UpdateState(gomock.Any(), gomock.Any(), filemodel.StateCompleted).
			DoAndReturn(func(ctx context.Context, id string, state filemodel.State) error {
				assert.Equal(s.T(), fileID, id)
				return updateErr
			})

		req := UploadRequest{
			UserID:      userID,
			Target:      "user_avatar",
			Name:        "avatar.png",
			ContentType: "image/png",
			Reader:      bytes.NewReader([]byte("test")),
		}

		resp, err := s.service.Upload(s.ctx, req)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "failed to update file state to completed"))
		assert.Empty(s.T(), resp.ID)
	})
}

func (s *ServiceTestSuite) TestGetURL() {
	s.Run("should get a file URL successfully", func() {
		fileID := uuid.New().String()
		expectedURL := "https://storage.example.com/user-avatar/avatar.png"
		expectedFile := filemodel.File{
			ID:        fileID,
			UserID:    uuid.New().String(),
			Path:      "user-avatar/avatar.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StateCompleted,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}

		s.mockPostgresRepo.EXPECT().
			Get(gomock.Any(), fileID).
			Return(expectedFile, nil)

		s.mockMinioRepo.EXPECT().
			GetURL(gomock.Any(), expectedFile.Path).
			Return(expectedURL, nil)

		url, err := s.service.GetURL(s.ctx, fileID)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedURL, url)
	})

	s.Run("should return error when file not found in database", func() {
		fileID := uuid.New().String()
		dbErr := errors.New("file not found")

		s.mockPostgresRepo.EXPECT().
			Get(gomock.Any(), fileID).
			Return(filemodel.File{}, dbErr)

		url, err := s.service.GetURL(s.ctx, fileID)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "failed to get file from database"))
		assert.Empty(s.T(), url)
	})

	s.Run("should return error when minio GetURL fails", func() {
		fileID := uuid.New().String()
		expectedFile := filemodel.File{
			ID:        fileID,
			UserID:    uuid.New().String(),
			Path:      "user-avatar/avatar.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StateCompleted,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}
		minioErr := errors.New("minio error")

		s.mockPostgresRepo.EXPECT().
			Get(gomock.Any(), fileID).
			Return(expectedFile, nil)

		s.mockMinioRepo.EXPECT().
			GetURL(gomock.Any(), expectedFile.Path).
			Return("", minioErr)

		url, err := s.service.GetURL(s.ctx, fileID)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "failed to get file URL from static storage"))
		assert.Empty(s.T(), url)
	})
}

func (s *ServiceTestSuite) TestDelete() {
	s.Run("should delete a file successfully", func() {
		fileID := uuid.New().String()
		expectedFile := filemodel.File{
			ID:        fileID,
			UserID:    uuid.New().String(),
			Path:      "user-avatar/avatar.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StateCompleted,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}

		s.mockPostgresRepo.EXPECT().
			Get(gomock.Any(), fileID).
			Return(expectedFile, nil)

		s.mockMinioRepo.EXPECT().
			Delete(gomock.Any(), expectedFile.Path).
			Return(nil)

		s.mockPostgresRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Return(nil)

		err := s.service.Delete(s.ctx, fileID)
		assert.NoError(s.T(), err)
	})

	s.Run("should return error when file not found in database", func() {
		fileID := uuid.New().String()
		dbErr := errors.New("file not found")

		s.mockPostgresRepo.EXPECT().
			Get(gomock.Any(), fileID).
			Return(filemodel.File{}, dbErr)

		err := s.service.Delete(s.ctx, fileID)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "failed to get file from database"))
	})

	s.Run("should return error when minio Delete fails", func() {
		fileID := uuid.New().String()
		expectedFile := filemodel.File{
			ID:        fileID,
			UserID:    uuid.New().String(),
			Path:      "user-avatar/avatar.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StateCompleted,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}
		minioErr := errors.New("minio delete error")

		s.mockPostgresRepo.EXPECT().
			Get(gomock.Any(), fileID).
			Return(expectedFile, nil)

		s.mockMinioRepo.EXPECT().
			Delete(gomock.Any(), expectedFile.Path).
			Return(minioErr)

		err := s.service.Delete(s.ctx, fileID)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "failed to delete file from static storage"))
	})

	s.Run("should return error when postgres Delete fails", func() {
		fileID := uuid.New().String()
		expectedFile := filemodel.File{
			ID:        fileID,
			UserID:    uuid.New().String(),
			Path:      "user-avatar/avatar.png",
			Size:      1024,
			Extension: filemodel.ExtensionPNG,
			State:     filemodel.StateCompleted,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}
		dbErr := errors.New("database delete error")

		s.mockPostgresRepo.EXPECT().
			Get(gomock.Any(), fileID).
			Return(expectedFile, nil)

		s.mockMinioRepo.EXPECT().
			Delete(gomock.Any(), expectedFile.Path).
			Return(nil)

		s.mockPostgresRepo.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Return(dbErr)

		err := s.service.Delete(s.ctx, fileID)
		assert.Error(s.T(), err)
		assert.True(s.T(), strings.Contains(err.Error(), "failed to delete file record from database"))
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

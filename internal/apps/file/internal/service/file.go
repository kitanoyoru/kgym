package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/minio"
	filemodel "github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/models/file"
	"github.com/kitanoyoru/kgym/internal/apps/file/internal/repository/postgres"
	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/file/pkg/validator"
)

type Service struct {
	cfg Config

	minioRepository    minio.IRepository
	postgresRepository postgres.IRepository
}

func New(cfg Config, minioRepository minio.IRepository, postgresRepository postgres.IRepository) *Service {
	return &Service{
		cfg:                cfg,
		minioRepository:    minioRepository,
		postgresRepository: postgresRepository,
	}
}

func (s *Service) Upload(ctx context.Context, req UploadRequest) (UploadResponse, error) {
	if err := pkgValidator.Validate.StructCtx(ctx, req); err != nil {
		return UploadResponse{}, err
	}

	bucket, ok := s.cfg.Buckets[req.Target]
	if !ok || bucket == "" {
		return UploadResponse{}, ErrBucketNotFound
	}

	extension, err := filemodel.ExtensionFromFileName(req.Name)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("invalid file extension: %w", err)
	}

	now := time.Now()
	fileID := uuid.New().String()
	path := filepath.Join(bucket, req.Name)

	file := filemodel.File{
		ID:        fileID,
		UserID:    req.UserID,
		Path:      path,
		Size:      0, // TODO: get size from reader
		Extension: extension,
		State:     filemodel.StatePending,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	if err := s.postgresRepository.Create(ctx, file); err != nil {
		return UploadResponse{}, fmt.Errorf("failed to create file record: %w", err)
	}

	minioReq := minio.UploadRequest{
		Bucket:      bucket,
		Name:        req.Name,
		ContentType: req.ContentType,
		Reader:      req.Reader,
	}

	if _, err := s.minioRepository.Upload(ctx, minioReq); err != nil {
		if updateErr := s.postgresRepository.UpdateState(ctx, fileID, filemodel.StateFailed); updateErr != nil {
			return UploadResponse{}, fmt.Errorf("upload failed: %w; state update failed: %w", err, updateErr)
		}

		return UploadResponse{}, fmt.Errorf("upload failed: %w", err)
	}

	if err := s.postgresRepository.UpdateState(ctx, fileID, filemodel.StateCompleted); err != nil {
		return UploadResponse{}, fmt.Errorf("failed to update file state to completed: %w", err)
	}

	return UploadResponse{
		ID: fileID,
	}, nil
}

func (s *Service) GetURL(ctx context.Context, id string) (string, error) {
	file, err := s.postgresRepository.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get file from database: %w", err)
	}

	url, err := s.minioRepository.GetURL(ctx, file.Path)
	if err != nil {
		return "", fmt.Errorf("failed to get file URL from static storage: %w", err)
	}

	return url, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	file, err := s.postgresRepository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get file from database: %w", err)
	}

	if err := s.minioRepository.Delete(ctx, file.Path); err != nil {
		return fmt.Errorf("failed to delete file from static storage: %w", err)
	}

	if err := s.postgresRepository.Delete(ctx, postgres.WithID(file.ID)); err != nil {
		return fmt.Errorf("failed to delete file record from database: %w", err)
	}

	return nil
}

package service

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

var (
	ErrBucketNotFound = errors.New("bucket not found")
	ErrFileNotFound   = errors.New("file not found")
)

type Config struct {
	Buckets map[string]string
}

type IService interface {
	Upload(ctx context.Context, req UploadRequest) (UploadResponse, error)
	GetURL(ctx context.Context, id string) (string, error)
	Delete(ctx context.Context, id string) error
}

type (
	UploadRequest struct {
		UserID      string    `validate:"required,uuid"`
		Target      string    `validate:"required,oneof=user_avatar"`
		Name        string    `validate:"required,min=1,max=255"`
		ContentType string    `validate:"required,min=1,max=255"`
		Reader      io.Reader `validate:"required"`
	}

	UploadResponse struct {
		ID string
	}
)

package minio

import (
	"context"
	"errors"
	"io"
	"time"
)

const (
	presignedURLExpiration     = 24 * time.Hour
	expectedPathPartsCount     = 2
	responseContentDisposition = "response-content-disposition"
)

var (
	ErrInvalidPath       = errors.New("invalid path: must be in format 'bucket/object'")
	ErrUnableToDetectExt = errors.New("unable to determine file extension")
)

type ConstructorOption func(*ConstructorOptions)

type ConstructorOptions struct {
	Buckets []string
}

func WithBuckets(buckets ...string) ConstructorOption {
	return func(o *ConstructorOptions) {
		o.Buckets = buckets
	}
}

type UploadRequest struct {
	Bucket, Name, ContentType string
	Reader                    io.Reader
}

type UploadResponse struct {
	URL string

	Extension string
	Size      int64
}

type IRepository interface {
	Upload(ctx context.Context, req UploadRequest) (UploadResponse, error)
	GetURL(ctx context.Context, path string) (string, error)
	Delete(ctx context.Context, path string) error
}

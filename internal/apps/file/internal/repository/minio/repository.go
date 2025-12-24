package minio

import (
	"context"
	"io"
)

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

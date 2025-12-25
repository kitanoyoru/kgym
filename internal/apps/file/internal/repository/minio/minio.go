package minio

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type Repository struct {
	minioClient *minio.Client
}

var _ IRepository = (*Repository)(nil)

func New(minioClient *minio.Client) *Repository {
	return &Repository{minioClient}
}

func (r *Repository) Upload(ctx context.Context, req UploadRequest) (UploadResponse, error) {
	name, err := r.getFileName(req.Name)
	if err != nil {
		return UploadResponse{}, err
	}

	uploadInfo, err := r.minioClient.PutObject(ctx, req.Bucket, name, req.Reader, -1, minio.PutObjectOptions{
		ContentType: req.ContentType,
	})
	if err != nil {
		return UploadResponse{}, err
	}

	extension := filepath.Ext(req.Name)
	if len(extension) > 0 {
		extension = strings.TrimPrefix(extension, ".")
	}

	path := fmt.Sprintf("%s/%s", req.Bucket, req.Name)

	return UploadResponse{
		URL:       r.getFileURL(path),
		Extension: extension,
		Size:      uploadInfo.Size,
	}, nil
}

func (r *Repository) GetURL(ctx context.Context, path string) (string, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return "", errors.New("invalid path")
	}

	bucket, object := parts[0], parts[1]

	params := make(url.Values)
	params.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", object))

	u, err := r.minioClient.PresignedGetObject(ctx, bucket, object, time.Hour*24, params)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func (r *Repository) Delete(ctx context.Context, path string) error {
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return errors.New("invalid path")
	}

	bucket, object := parts[0], parts[1]

	return r.minioClient.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
}

func (r *Repository) getFileName(src string) (string, error) {
	index := strings.LastIndex(src, ".")
	if index < 0 {
		return "", errors.New("unable to determine file type")
	}

	name := strings.Replace(uuid.New().String(), "-", "", -1)
	name += src[index:]

	return name, nil
}

func (r *Repository) getFileURL(path string) string {
	return fmt.Sprintf("%s/%s", r.minioClient.EndpointURL(), path)
}

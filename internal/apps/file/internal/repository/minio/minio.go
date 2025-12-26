package minio

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/errgroup"
)

type Repository struct {
	minioClient *minio.Client
}

var _ IRepository = (*Repository)(nil)

func New(minioClient *minio.Client, options ...ConstructorOption) (*Repository, error) {
	var opts ConstructorOptions
	for _, option := range options {
		option(&opts)
	}

	repository := &Repository{minioClient: minioClient}

	if len(opts.Buckets) > 0 {
		if err := repository.makeBuckets(context.Background(), opts.Buckets...); err != nil {
			return nil, err
		}
	}

	return repository, nil
}

func (r *Repository) Upload(ctx context.Context, req UploadRequest) (UploadResponse, error) {
	objectName, err := r.generateObjectName(req.Name)
	if err != nil {
		return UploadResponse{}, err
	}

	uploadInfo, err := r.minioClient.PutObject(ctx, req.Bucket, objectName, req.Reader, -1, minio.PutObjectOptions{
		ContentType: req.ContentType,
	})
	if err != nil {
		return UploadResponse{}, err
	}

	extension := r.extractExtension(uploadInfo.Key, req.Name)
	filePath := r.buildPath(req.Bucket, uploadInfo.Key)
	fileURL := r.buildFileURL(filePath)

	return UploadResponse{
		URL:       fileURL,
		Extension: extension,
		Size:      uploadInfo.Size,
	}, nil
}

func (r *Repository) GetURL(ctx context.Context, path string) (string, error) {
	bucket, object, err := r.parsePath(path)
	if err != nil {
		return "", err
	}

	params := make(url.Values)
	params.Set(responseContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", object))

	presignedURL, err := r.minioClient.PresignedGetObject(ctx, bucket, object, presignedURLExpiration, params)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (r *Repository) Delete(ctx context.Context, path string) error {
	bucket, object, err := r.parsePath(path)
	if err != nil {
		return err
	}

	return r.minioClient.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
}

func (r *Repository) generateObjectName(fileName string) (string, error) {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return "", ErrUnableToDetectExt
	}

	uniqueID := strings.ReplaceAll(uuid.New().String(), "-", "")
	return uniqueID + ext, nil
}

func (r *Repository) extractExtension(objectKey, originalFileName string) string {
	ext := filepath.Ext(objectKey)
	if ext != "" {
		return strings.TrimPrefix(ext, ".")
	}

	ext = filepath.Ext(originalFileName)
	if ext != "" {
		return strings.TrimPrefix(ext, ".")
	}

	return ""
}

func (r *Repository) buildPath(bucket, objectKey string) string {
	return bucket + "/" + objectKey
}

func (r *Repository) buildFileURL(path string) string {
	endpoint := r.minioClient.EndpointURL()
	if endpoint == nil {
		return ""
	}

	return endpoint.String() + "/" + path
}

func (r *Repository) parsePath(path string) (string, string, error) {
	if path == "" {
		return "", "", ErrInvalidPath
	}

	path = strings.Trim(path, "/")
	parts := strings.SplitN(path, "/", expectedPathPartsCount)

	if len(parts) != expectedPathPartsCount {
		return "", "", ErrInvalidPath
	}

	bucket := strings.TrimSpace(parts[0])
	object := strings.TrimSpace(parts[1])

	return bucket, object, nil
}

func (r *Repository) makeBuckets(ctx context.Context, buckets ...string) error {
	if len(buckets) == 0 {
		return nil
	}

	var errGroup errgroup.Group

	for i := range buckets {
		bucket := buckets[i]
		errGroup.Go(func() error {
			return r.minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		})
	}

	return errGroup.Wait()
}

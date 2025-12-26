package decorators

import (
	"context"

	"github.com/kitanoyoru/kgym/internal/apps/file/internal/service"
	pkgLogger "github.com/kitanoyoru/kgym/pkg/logger"
	"go.uber.org/zap"
)

type loggingDecorator struct {
	svc service.IService
}

func Logging(svc service.IService) service.IService {
	return &loggingDecorator{svc: svc}
}

func (ld *loggingDecorator) Upload(ctx context.Context, req service.UploadRequest) (service.UploadResponse, error) {
	logger, err := pkgLogger.FromContext(ctx)
	if err != nil {
		return service.UploadResponse{}, err
	}

	logger.With(zap.String("service", "file")).Info("uploading file",
		zap.String("user_id", req.UserID),
		zap.String("target", req.Target),
		zap.String("name", req.Name),
		zap.String("content_type", req.ContentType),
	)

	return ld.svc.Upload(ctx, req)
}

func (ld *loggingDecorator) GetURL(ctx context.Context, id string) (string, error) {
	logger, err := pkgLogger.FromContext(ctx)
	if err != nil {
		return "", err
	}

	logger.With(zap.String("service", "file")).Info("getting file URL",
		zap.String("file_id", id),
	)

	return ld.svc.GetURL(ctx, id)

}

func (ld *loggingDecorator) Delete(ctx context.Context, id string) error {
	logger, err := pkgLogger.FromContext(ctx)
	if err != nil {
		return err
	}

	logger.With(zap.String("service", "file")).Info("deleting file",
		zap.String("file_id", id),
	)

	return ld.svc.Delete(ctx, id)

}

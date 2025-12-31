package decorators

import (
	"context"
	"log"

	"github.com/kitanoyoru/kgym/internal/apps/file/internal/service"
)

type loggingDecorator struct {
	svc service.IService
}

func Logging(svc service.IService) service.IService {
	return &loggingDecorator{svc: svc}
}

func (ld *loggingDecorator) Upload(ctx context.Context, req service.UploadRequest) (service.UploadResponse, error) {
	log.Printf("uploading file: user_id=%s target=%s name=%s content_type=%s",
		req.UserID,
		req.Target,
		req.Name,
		req.ContentType,
	)

	return ld.svc.Upload(ctx, req)
}

func (ld *loggingDecorator) GetURL(ctx context.Context, id string) (string, error) {
	log.Printf("getting file URL: file_id=%s", id)

	return ld.svc.GetURL(ctx, id)
}

func (ld *loggingDecorator) Delete(ctx context.Context, id string) error {
	log.Printf("deleting file: file_id=%s", id)

	return ld.svc.Delete(ctx, id)
}

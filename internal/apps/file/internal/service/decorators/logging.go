package decorators

import (
	"context"

	"github.com/kitanoyoru/kgym/internal/apps/file/internal/service"
	"github.com/rs/zerolog/log"
)

type loggingDecorator struct {
	svc service.IService
}

func Logging(svc service.IService) service.IService {
	return &loggingDecorator{svc: svc}
}

func (ld *loggingDecorator) Upload(ctx context.Context, req service.UploadRequest) (service.UploadResponse, error) {
	log.Info().
		Str("user_id", req.UserID).
		Str("target", req.Target).
		Str("name", req.Name).
		Str("content_type", req.ContentType).
		Msg("uploading file")

	return ld.svc.Upload(ctx, req)
}

func (ld *loggingDecorator) GetURL(ctx context.Context, id string) (string, error) {
	log.Info().
		Str("file_id", id).
		Msg("getting file URL")

	return ld.svc.GetURL(ctx, id)
}

func (ld *loggingDecorator) Delete(ctx context.Context, id string) error {
	log.Info().
		Str("file_id", id).
		Msg("deleting file")

	return ld.svc.Delete(ctx, id)
}

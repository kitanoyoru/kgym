package decorators

import (
	"context"

	"github.com/google/uuid"
	"github.com/kitanoyoru/kgym/internal/apps/file/internal/service"
	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/file/pkg/validator"
)

type validateDecorator struct {
	svc service.IService
}

func Validate(svc service.IService) service.IService {
	return &validateDecorator{svc: svc}
}

func (vd *validateDecorator) Upload(ctx context.Context, req service.UploadRequest) (service.UploadResponse, error) {
	if err := pkgValidator.Validate.StructCtx(ctx, req); err != nil {
		return service.UploadResponse{}, err
	}

	return vd.svc.Upload(ctx, req)
}

func (vd *validateDecorator) GetURL(ctx context.Context, id string) (string, error) {
	if err := uuid.Validate(id); err != nil {
		return "", err
	}

	return vd.svc.GetURL(ctx, id)
}

func (vd *validateDecorator) Delete(ctx context.Context, id string) error {
	if err := uuid.Validate(id); err != nil {
		return err
	}

	return vd.svc.Delete(ctx, id)
}

package decorators

import (
	"context"

	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/service"
	pkgValidator "github.com/kitanoyoru/kgym/internal/apps/sso/pkg/validator"
)

type validateDecorator struct {
	svc service.IService
}

func Validate(svc service.IService) service.IService {
	return &validateDecorator{svc: svc}
}

func (vd *validateDecorator) CreateToken(ctx context.Context, req service.CreateTokenRequest) (string, error) {
	if err := pkgValidator.Validate.StructCtx(ctx, req); err != nil {
		return "", err
	}

	return vd.svc.CreateToken(ctx, req)
}

func (vd *validateDecorator) ListTokens(ctx context.Context, filters ...service.Filter) ([]token.Token, error) {
	return vd.svc.ListTokens(ctx, filters...)
}

func (vd *validateDecorator) UpdateToken(ctx context.Context, id, token string) error {
	if err := pkgValidator.Validate.VarCtx(ctx, id, "required,uuid"); err != nil {
		return err
	}

	if err := pkgValidator.Validate.VarCtx(ctx, token, "required,min=1,max=255"); err != nil {
		return err
	}

	return vd.svc.UpdateToken(ctx, id, token)
}

func (vd *validateDecorator) DeleteToken(ctx context.Context, id string) error {
	if err := pkgValidator.Validate.VarCtx(ctx, id, "required,uuid"); err != nil {
		return err
	}

	return vd.svc.DeleteToken(ctx, id)
}

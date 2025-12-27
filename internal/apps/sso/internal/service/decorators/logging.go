package decorators

import (
	"context"

	"go.uber.org/zap"

	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/service"
	pkgLogger "github.com/kitanoyoru/kgym/pkg/logger"
)

type loggingDecorator struct {
	svc service.IService
}

func Logging(svc service.IService) service.IService {
	return &loggingDecorator{svc: svc}
}

func (ld *loggingDecorator) CreateToken(ctx context.Context, req service.CreateTokenRequest) (string, error) {
	logger, err := pkgLogger.FromContext(ctx)
	if err != nil {
		return "", err
	}

	logger.With(zap.String("service", "sso")).Info("creating token",
		zap.String("user_id", req.UserID),
		zap.String("token_type", req.TokenType.String()),
		zap.String("token", req.Token),
	)

	return ld.svc.CreateToken(ctx, req)
}

func (ld *loggingDecorator) ListTokens(ctx context.Context, filters ...service.Filter) ([]token.Token, error) {
	logger, err := pkgLogger.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	logger.With(zap.String("service", "sso")).Info("listing tokens",
		zap.Any("filters", filters),
	)

	return ld.svc.ListTokens(ctx, filters...)
}

func (ld *loggingDecorator) UpdateToken(ctx context.Context, id, token string) error {
	logger, err := pkgLogger.FromContext(ctx)
	if err != nil {
		return err
	}

	logger.With(zap.String("service", "sso")).Info("updating token",
		zap.String("id", id),
		zap.String("token", token),
	)
	return ld.svc.UpdateToken(ctx, id, token)
}

func (ld *loggingDecorator) DeleteToken(ctx context.Context, id string) error {
	logger, err := pkgLogger.FromContext(ctx)
	if err != nil {
		return err
	}

	logger.With(zap.String("service", "sso")).Info("deleting token",
		zap.String("id", id),
	)

	return ld.svc.DeleteToken(ctx, id)
}

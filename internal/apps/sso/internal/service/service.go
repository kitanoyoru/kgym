package service

import (
	"context"

	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
)

type IService interface {
	CreateToken(ctx context.Context, req CreateTokenRequest) (string, error)
	ListTokens(ctx context.Context, filters ...Filter) ([]token.Token, error)
	UpdateToken(ctx context.Context, id, token string) error
	DeleteToken(ctx context.Context, id string) error
}

type (
	CreateTokenRequest struct {
		UserID    string
		TokenType token.TokenType
		Token     string
	}
)

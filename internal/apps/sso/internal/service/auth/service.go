package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

const (
	IssuerServiceName = "sso.kgym"

	AccessTokenTTL  = 15 * time.Minute   // 15 minutes
	RefreshTokenTTL = 7 * 24 * time.Hour // 7 days
)

type IService interface {
	PasswordGrant(ctx context.Context, req PasswordGrantRequest) (PasswordGrantResponse, error)
	RefreshTokenGrant(ctx context.Context, req RefreshTokenGrantRequest) (RefreshTokenGrantResponse, error)
}

type (
	PasswordGrantRequest struct {
		Email    string
		Password string

		ClientID string
	}

	PasswordGrantResponse struct {
		AccessToken  string
		RefreshToken string
	}
)

type (
	RefreshTokenGrantRequest struct {
		RefreshToken string
	}

	RefreshTokenGrantResponse struct {
		AccessToken  string
		RefreshToken string
	}
)

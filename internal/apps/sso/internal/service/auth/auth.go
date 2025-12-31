package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	keyrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/key"
	tokenrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token"
	userrepo "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user"
)

var _ IService = (*Service)(nil)

type Service struct {
	userRepository  userrepo.IRepository
	tokenRepository tokenrepo.IRepository
	keyRepository   keyrepo.IRepository
}

func NewService(userRepository userrepo.IRepository, tokenRepository tokenrepo.IRepository, keyRepository keyrepo.IRepository) *Service {
	return &Service{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		keyRepository:   keyRepository,
	}
}

func (s *Service) PasswordGrant(ctx context.Context, req PasswordGrantRequest) (PasswordGrantResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return PasswordGrantResponse{}, err
	}

	if user.Password != req.Password {
		return PasswordGrantResponse{}, ErrInvalidCredentials
	}

	accessToken, err := s.issueAccessToken(ctx, user.ID, req.ClientID)
	if err != nil {
		return PasswordGrantResponse{}, err
	}

	refreshToken, err := s.issueRefreshToken(ctx, user.ID, req.ClientID)
	if err != nil {
		return PasswordGrantResponse{}, err
	}

	return PasswordGrantResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshTokenGrant(ctx context.Context, req RefreshTokenGrantRequest) (RefreshTokenGrantResponse, error) {
	token, err := s.tokenRepository.GetByTokenHash(ctx, req.RefreshToken)
	if err != nil {
		return RefreshTokenGrantResponse{}, err
	}

	if token.Revoked || token.ExpiresAt.Before(time.Now()) {
		return RefreshTokenGrantResponse{}, ErrInvalidRefreshToken
	}

	err = s.tokenRepository.Revoke(ctx, req.RefreshToken)
	if err != nil {
		return RefreshTokenGrantResponse{}, err
	}

	access, err := s.issueAccessToken(ctx, token.Subject, token.ClientID)
	if err != nil {
		return RefreshTokenGrantResponse{}, err
	}

	refresh, err := s.issueRefreshToken(ctx, token.Subject, token.ClientID)
	if err != nil {
		return RefreshTokenGrantResponse{}, err
	}

	return RefreshTokenGrantResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) issueAccessToken(ctx context.Context, subject, clientID string) (string, error) {
	key, err := s.keyRepository.GetCurrentSigningKey(ctx)
	if err != nil {
		return "", err
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		Subject:   subject,
		Audience:  jwt.ClaimStrings{clientID},
		Issuer:    IssuerServiceName,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
	})

	return token.SignedString(key.Private)
}

func (s *Service) issueRefreshToken(ctx context.Context, subject, clientID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	sum := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(sum[:])

	err := s.tokenRepository.Create(ctx, tokenentity.Token{
		ID:        uuid.New().String(),
		Subject:   subject,
		ClientID:  clientID,
		TokenType: tokenentity.TypeRefresh,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

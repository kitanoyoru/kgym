package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	keyentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/key"
	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	keymocks "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/key/mocks"
	tokenmocks "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token/mocks"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token/models/token"
	usermocks "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user/mocks"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user/models"
)

func TestService_PasswordGrant(t *testing.T) {
	t.Run("should grant tokens successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		email := "test@example.com"
		password := "password123"
		clientID := "client-123"
		userID := "user-123"

		user := usermodel.User{
			ID:       userID,
			Email:    email,
			Password: password,
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		key := keyentity.Key{
			ID:        "key-123",
			Private:   privateKey,
			Public:    privateKey.Public(),
			Algorithm: "RS256",
			Active:    true,
		}

		userRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		keyRepo.EXPECT().
			GetCurrentSigningKey(ctx).
			Return(key, nil)

		tokenRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, token tokenentity.Token) error {
				assert.Equal(t, userID, token.Subject)
				assert.Equal(t, clientID, token.ClientID)
				assert.Equal(t, tokenentity.TokenTypeRefresh, token.TokenType)
				assert.False(t, token.Revoked)
				assert.NotEmpty(t, token.TokenHash)
				return nil
			})

		req := PasswordGrantRequest{
			Email:    email,
			Password: password,
			ClientID: clientID,
		}

		resp, err := service.PasswordGrant(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)

		parser := jwt.NewParser()
		claims := jwt.RegisteredClaims{}
		_, err = parser.ParseWithClaims(resp.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
			return privateKey.Public(), nil
		})
		require.NoError(t, err)
		assert.Equal(t, userID, claims.Subject)
		assert.Equal(t, []string{clientID}, []string(claims.Audience))
		assert.Equal(t, IssuerServiceName, claims.Issuer)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		email := "notfound@example.com"

		userRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(usermodel.User{}, errors.New("user not found"))

		req := PasswordGrantRequest{
			Email:    email,
			Password: "password",
			ClientID: "client-123",
		}

		resp, err := service.PasswordGrant(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		email := "test@example.com"
		password := "correct-password"
		wrongPassword := "wrong-password"

		user := usermodel.User{
			ID:       "user-123",
			Email:    email,
			Password: password,
		}

		userRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		req := PasswordGrantRequest{
			Email:    email,
			Password: wrongPassword,
			ClientID: "client-123",
		}

		resp, err := service.PasswordGrant(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCredentials, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("should return error when key repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		email := "test@example.com"
		password := "password123"

		user := usermodel.User{
			ID:       "user-123",
			Email:    email,
			Password: password,
		}

		userRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		keyRepo.EXPECT().
			GetCurrentSigningKey(ctx).
			Return(keyentity.Key{}, errors.New("key repository error"))

		req := PasswordGrantRequest{
			Email:    email,
			Password: password,
			ClientID: "client-123",
		}

		resp, err := service.PasswordGrant(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("should return error when token repository create fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		email := "test@example.com"
		password := "password123"

		user := usermodel.User{
			ID:       "user-123",
			Email:    email,
			Password: password,
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		key := keyentity.Key{
			ID:        "key-123",
			Private:   privateKey,
			Public:    privateKey.Public(),
			Algorithm: "RS256",
			Active:    true,
		}

		userRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		keyRepo.EXPECT().
			GetCurrentSigningKey(ctx).
			Return(key, nil)

		tokenRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("token repository error"))

		req := PasswordGrantRequest{
			Email:    email,
			Password: password,
			ClientID: "client-123",
		}

		resp, err := service.PasswordGrant(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})
}

func TestService_RefreshTokenGrant(t *testing.T) {
	t.Run("should grant new tokens successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		refreshTokenHash := "token-hash-123"
		subject := "user-123"
		clientID := "client-123"

		token := tokenmodel.Token{
			ID:        "token-123",
			Subject:   subject,
			ClientID:  clientID,
			TokenHash: refreshTokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   false,
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		key := keyentity.Key{
			ID:        "key-123",
			Private:   privateKey,
			Public:    privateKey.Public(),
			Algorithm: "RS256",
			Active:    true,
		}

		tokenRepo.EXPECT().
			GetByTokenHash(ctx, refreshTokenHash).
			Return(token, nil)

		tokenRepo.EXPECT().
			Revoke(ctx, refreshTokenHash).
			Return(nil)

		keyRepo.EXPECT().
			GetCurrentSigningKey(ctx).
			Return(key, nil)

		tokenRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, newToken tokenentity.Token) error {
				assert.Equal(t, subject, newToken.Subject)
				assert.Equal(t, clientID, newToken.ClientID)
				assert.Equal(t, tokenentity.TokenTypeRefresh, newToken.TokenType)
				assert.False(t, newToken.Revoked)
				return nil
			})

		req := RefreshTokenGrantRequest{
			RefreshToken: refreshTokenHash,
		}

		resp, err := service.RefreshTokenGrant(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)

		parser := jwt.NewParser()
		claims := jwt.RegisteredClaims{}
		_, err = parser.ParseWithClaims(resp.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
			return privateKey.Public(), nil
		})
		require.NoError(t, err)
		assert.Equal(t, subject, claims.Subject)
		assert.Equal(t, []string{clientID}, []string(claims.Audience))
	})

	t.Run("should return error when token not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		refreshTokenHash := "invalid-hash"

		tokenRepo.EXPECT().
			GetByTokenHash(ctx, refreshTokenHash).
			Return(tokenmodel.Token{}, errors.New("token not found"))

		req := RefreshTokenGrantRequest{
			RefreshToken: refreshTokenHash,
		}

		resp, err := service.RefreshTokenGrant(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("should return error when token is revoked", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		refreshTokenHash := "token-hash-123"

		token := tokenmodel.Token{
			ID:        "token-123",
			Subject:   "user-123",
			ClientID:  "client-123",
			TokenHash: refreshTokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   true,
		}

		tokenRepo.EXPECT().
			GetByTokenHash(ctx, refreshTokenHash).
			Return(token, nil)

		req := RefreshTokenGrantRequest{
			RefreshToken: refreshTokenHash,
		}

		resp, err := service.RefreshTokenGrant(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRefreshToken, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("should return error when token is expired", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		refreshTokenHash := "token-hash-123"

		token := tokenmodel.Token{
			ID:        "token-123",
			Subject:   "user-123",
			ClientID:  "client-123",
			TokenHash: refreshTokenHash,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Revoked:   false,
		}

		tokenRepo.EXPECT().
			GetByTokenHash(ctx, refreshTokenHash).
			Return(token, nil)

		req := RefreshTokenGrantRequest{
			RefreshToken: refreshTokenHash,
		}

		resp, err := service.RefreshTokenGrant(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRefreshToken, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("should return error when revoke fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		refreshTokenHash := "token-hash-123"

		token := tokenmodel.Token{
			ID:        "token-123",
			Subject:   "user-123",
			ClientID:  "client-123",
			TokenHash: refreshTokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   false,
		}

		tokenRepo.EXPECT().
			GetByTokenHash(ctx, refreshTokenHash).
			Return(token, nil)

		tokenRepo.EXPECT().
			Revoke(ctx, refreshTokenHash).
			Return(errors.New("revoke error"))

		req := RefreshTokenGrantRequest{
			RefreshToken: refreshTokenHash,
		}

		resp, err := service.RefreshTokenGrant(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("should return error when key repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		refreshTokenHash := "token-hash-123"

		token := tokenmodel.Token{
			ID:        "token-123",
			Subject:   "user-123",
			ClientID:  "client-123",
			TokenHash: refreshTokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   false,
		}

		tokenRepo.EXPECT().
			GetByTokenHash(ctx, refreshTokenHash).
			Return(token, nil)

		tokenRepo.EXPECT().
			Revoke(ctx, refreshTokenHash).
			Return(nil)

		keyRepo.EXPECT().
			GetCurrentSigningKey(ctx).
			Return(keyentity.Key{}, errors.New("key repository error"))

		req := RefreshTokenGrantRequest{
			RefreshToken: refreshTokenHash,
		}

		resp, err := service.RefreshTokenGrant(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})

	t.Run("should return error when token repository create fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := usermocks.NewMockIRepository(ctrl)
		tokenRepo := tokenmocks.NewMockIRepository(ctrl)
		keyRepo := keymocks.NewMockIRepository(ctrl)

		service := &Service{
			userRepository:  userRepo,
			tokenRepository: tokenRepo,
			keyRepository:   keyRepo,
		}

		ctx := context.Background()
		refreshTokenHash := "token-hash-123"

		token := tokenmodel.Token{
			ID:        "token-123",
			Subject:   "user-123",
			ClientID:  "client-123",
			TokenHash: refreshTokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   false,
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		key := keyentity.Key{
			ID:        "key-123",
			Private:   privateKey,
			Public:    privateKey.Public(),
			Algorithm: "RS256",
			Active:    true,
		}

		tokenRepo.EXPECT().
			GetByTokenHash(ctx, refreshTokenHash).
			Return(token, nil)

		tokenRepo.EXPECT().
			Revoke(ctx, refreshTokenHash).
			Return(nil)

		keyRepo.EXPECT().
			GetCurrentSigningKey(ctx).
			Return(key, nil)

		tokenRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("token repository error"))

		req := RefreshTokenGrantRequest{
			RefreshToken: refreshTokenHash,
		}

		resp, err := service.RefreshTokenGrant(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, resp.AccessToken)
		assert.Empty(t, resp.RefreshToken)
	})
}

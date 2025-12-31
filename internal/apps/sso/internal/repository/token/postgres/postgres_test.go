package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	tokenentity "github.com/kitanoyoru/kgym/internal/apps/sso/internal/entity/token"
	tokenmodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token/models/token"
	"github.com/kitanoyoru/kgym/internal/apps/sso/migrations"
	"github.com/kitanoyoru/kgym/pkg/database/postgres"
	"github.com/kitanoyoru/kgym/pkg/testing/integration/cockroachdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type RepositoryTestSuite struct {
	suite.Suite

	db        *pgxpool.Pool
	container *cockroachdb.CockroachDBContainer
}

func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	container, err := cockroachdb.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup test container")

	s.container = container

	s.db, err = postgres.New(ctx, postgres.Config{
		URI: container.URI,
	})
	require.NoError(s.T(), err, "failed to create postgres client")

	err = migrations.Up(ctx, "pgx", container.URI)
	require.NoError(s.T(), err, "failed to run migrations")
}

func (s *RepositoryTestSuite) TearDownSuite() {
	if s.container != nil {
		_ = s.container.Terminate(s.T().Context())
	}
	if s.db != nil {
		s.db.Close()
	}
}

func (s *RepositoryTestSuite) SetupTest() {
	ctx := context.Background()
	_, err := s.db.Exec(ctx, "DELETE FROM tokens")
	require.NoError(s.T(), err, "failed to clean tokens table")
}

func (s *RepositoryTestSuite) TearDownTest() {
	ctx := context.Background()
	_, err := s.db.Exec(ctx, "DELETE FROM tokens")
	require.NoError(s.T(), err, "failed to clean tokens table")
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *RepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should create a token successfully", func() {
		tokenValue := "refresh-token-123"
		tokenHash := hashToken(tokenValue)
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		assert.NoError(s.T(), err)

		retrievedToken, err := repository.GetByTokenHash(ctx, tokenHash)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), token.ID, retrievedToken.ID)
		assert.Equal(s.T(), token.Subject, retrievedToken.Subject)
		assert.Equal(s.T(), token.ClientID, retrievedToken.ClientID)
		assert.Equal(s.T(), tokenmodel.TypeRefresh, retrievedToken.TokenType)
		assert.Equal(s.T(), token.TokenHash, retrievedToken.TokenHash)
		assert.Equal(s.T(), token.Revoked, retrievedToken.Revoked)
	})

	s.Run("should create multiple tokens with same subject and client_id", func() {
		subject := uuid.New().String()
		clientID := uuid.New().String()
		tokenHash1 := hashToken("refresh-token-1")
		token1 := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   subject,
			ClientID:  clientID,
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash1,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token1)
		require.NoError(s.T(), err)

		tokenHash2 := hashToken("refresh-token-2")
		token2 := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   subject,
			ClientID:  clientID,
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash2,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err = repository.Create(ctx, token2)
		assert.NoError(s.T(), err)

		retrievedToken1, err := repository.GetByTokenHash(ctx, tokenHash1)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), token1.ID, retrievedToken1.ID)

		retrievedToken2, err := repository.GetByTokenHash(ctx, tokenHash2)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), token2.ID, retrievedToken2.ID)
	})

	s.Run("should not create a token because of non-unique token_hash", func() {
		tokenValue := "duplicate-token-value"
		tokenHash := hashToken(tokenValue)
		token1 := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token1)
		require.NoError(s.T(), err)

		token2 := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err = repository.Create(ctx, token2)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a token because of empty subject", func() {
		tokenHash := hashToken("refresh-token-123")
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   "",
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		assert.Error(s.T(), err)
	})

	s.Run("should not create a token because of empty client_id", func() {
		tokenHash := hashToken("refresh-token-123")
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  "",
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		assert.Error(s.T(), err)
	})

	s.Run("should create a token with empty token_hash", func() {
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: "",
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		assert.NoError(s.T(), err)
	})

	s.Run("should create a token with non-uuid subject format", func() {
		tokenHash := hashToken("refresh-token-invalid-subject-" + uuid.New().String())
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   "invalid-uuid-format",
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		assert.NoError(s.T(), err)
	})
}

func (s *RepositoryTestSuite) TestGetByTokenHash() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should get a token by token hash successfully", func() {
		tokenValue := "refresh-token-by-hash"
		tokenHash := hashToken(tokenValue)
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		require.NoError(s.T(), err)

		retrievedToken, err := repository.GetByTokenHash(ctx, tokenHash)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), token.ID, retrievedToken.ID)
		assert.Equal(s.T(), token.Subject, retrievedToken.Subject)
		assert.Equal(s.T(), token.ClientID, retrievedToken.ClientID)
		assert.Equal(s.T(), token.TokenHash, retrievedToken.TokenHash)
		assert.Equal(s.T(), token.Revoked, retrievedToken.Revoked)
	})

	s.Run("should return error when token hash not found", func() {
		nonExistentHash := hashToken("non-existent-token")
		_, err := repository.GetByTokenHash(ctx, nonExistentHash)
		assert.Error(s.T(), err)
	})

	s.Run("should not get a token because it is deleted", func() {
		tokenValue := "refresh-token-to-delete"
		tokenHash := hashToken(tokenValue)
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		require.NoError(s.T(), err)

		_, err = s.db.Exec(ctx, "UPDATE tokens SET deleted_at = now() WHERE id = $1", token.ID)
		require.NoError(s.T(), err)

		_, err = repository.GetByTokenHash(ctx, tokenHash)
		assert.Error(s.T(), err)
	})
}

func (s *RepositoryTestSuite) TestRevoke() {
	ctx := context.Background()
	repository := New(s.db)

	s.Run("should revoke a token successfully", func() {
		tokenValue := "refresh-token-to-revoke"
		tokenHash := hashToken(tokenValue)
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		require.NoError(s.T(), err)

		retrievedToken, err := repository.GetByTokenHash(ctx, tokenHash)
		require.NoError(s.T(), err)
		assert.False(s.T(), retrievedToken.Revoked)

		err = repository.Revoke(ctx, tokenHash)
		assert.NoError(s.T(), err)

		retrievedToken, err = repository.GetByTokenHash(ctx, tokenHash)
		require.NoError(s.T(), err)
		assert.True(s.T(), retrievedToken.Revoked)
	})

	s.Run("should revoke a token that is already revoked", func() {
		tokenValue := "refresh-token-already-revoked"
		tokenHash := hashToken(tokenValue)
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   true,
		}

		err := repository.Create(ctx, token)
		require.NoError(s.T(), err)

		err = repository.Revoke(ctx, tokenHash)
		assert.NoError(s.T(), err)

		retrievedToken, err := repository.GetByTokenHash(ctx, tokenHash)
		require.NoError(s.T(), err)
		assert.True(s.T(), retrievedToken.Revoked)
	})

	s.Run("should return error when revoking non-existent token", func() {
		nonExistentHash := hashToken("non-existent-token")
		err := repository.Revoke(ctx, nonExistentHash)
		assert.NoError(s.T(), err)
	})

	s.Run("should not revoke a token because it is deleted", func() {
		tokenValue := "refresh-token-deleted"
		tokenHash := hashToken(tokenValue)
		token := tokenentity.Token{
			ID:        uuid.New().String(),
			Subject:   uuid.New().String(),
			ClientID:  uuid.New().String(),
			TokenType: tokenentity.TypeRefresh,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}

		err := repository.Create(ctx, token)
		require.NoError(s.T(), err)

		_, err = s.db.Exec(ctx, "UPDATE tokens SET deleted_at = now() WHERE id = $1", token.ID)
		require.NoError(s.T(), err)

		err = repository.Revoke(ctx, tokenHash)
		assert.NoError(s.T(), err)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

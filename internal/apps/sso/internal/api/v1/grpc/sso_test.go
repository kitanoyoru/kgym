package grpc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/sso/v1"
	"github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/key/redis"
	tokenpostgres "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/token/postgres"
	usermocks "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user/mocks"
	usermodel "github.com/kitanoyoru/kgym/internal/apps/sso/internal/repository/user/models"
	authservice "github.com/kitanoyoru/kgym/internal/apps/sso/internal/service/auth"
	keyservice "github.com/kitanoyoru/kgym/internal/apps/sso/internal/service/key"
	"github.com/kitanoyoru/kgym/internal/apps/sso/migrations"
	postgresdb "github.com/kitanoyoru/kgym/pkg/database/postgres"
	"github.com/kitanoyoru/kgym/pkg/testing/integration/cockroachdb"
	rediscontainer "github.com/kitanoyoru/kgym/pkg/testing/integration/redis"
	redisclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SSOServiceTestSuite struct {
	suite.Suite

	db                 *pgxpool.Pool
	rdb                *redisclient.Client
	cockroachContainer *cockroachdb.CockroachDBContainer
	redisContainer     *rediscontainer.RedisContainer
	ctrl               *gomock.Controller
	userRepo           *usermocks.MockIRepository
	ssoServer          *grpc.Server
	ssoClient          pb.SSOServiceClient
	ssoConn            *grpc.ClientConn
}

func (s *SSOServiceTestSuite) SetupSuite() {
	ctx := context.Background()

	cockroachContainer, err := cockroachdb.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup cockroachdb container")
	s.cockroachContainer = cockroachContainer

	s.db, err = postgresdb.New(ctx, postgresdb.Config{
		URI: cockroachContainer.URI,
	})
	require.NoError(s.T(), err, "failed to create postgres client")

	err = migrations.Up(ctx, "pgx", cockroachContainer.URI)
	require.NoError(s.T(), err, "failed to run migrations")

	redisContainer, err := rediscontainer.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup redis container")
	s.redisContainer = redisContainer

	parsedURI, err := url.Parse(redisContainer.URI)
	require.NoError(s.T(), err, "failed to parse Redis URI")

	address := strings.TrimPrefix(redisContainer.URI, "redis://")
	if parsedURI.Host != "" {
		address = parsedURI.Host
	}

	s.rdb = redisclient.NewClient(&redisclient.Options{
		Addr: address,
	})

	s.ctrl = gomock.NewController(s.T())
	s.userRepo = usermocks.NewMockIRepository(s.ctrl)

	tokenRepo := tokenpostgres.New(s.db)
	keyRepo, err := redis.New(ctx, s.rdb)
	require.NoError(s.T(), err, "failed to create key repository")

	authService := authservice.NewService(s.userRepo, tokenRepo, keyRepo)

	keyService := keyservice.NewService(keyRepo)

	ssoServer, err := NewSSOServer(authService, keyService)
	require.NoError(s.T(), err, "failed to create SSO server")

	s.ssoServer = grpc.NewServer()
	pb.RegisterSSOServiceServer(s.ssoServer, ssoServer)

	ssoListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err, "failed to create SSO listener")

	go func() {
		_ = s.ssoServer.Serve(ssoListener)
	}()

	time.Sleep(100 * time.Millisecond)

	ssoConn, err := grpc.NewClient(
		ssoListener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(s.T(), err, "failed to create SSO client")

	s.ssoConn = ssoConn
	s.ssoClient = pb.NewSSOServiceClient(ssoConn)
}

func (s *SSOServiceTestSuite) TearDownSuite() {
	if s.ctrl != nil {
		s.ctrl.Finish()
	}
	if s.ssoConn != nil {
		_ = s.ssoConn.Close()
	}
	if s.ssoServer != nil {
		s.ssoServer.GracefulStop()
	}
	if s.rdb != nil {
		_ = s.rdb.Close()
	}
	if s.db != nil {
		s.db.Close()
	}
	if s.redisContainer != nil {
		_ = s.redisContainer.Terminate(s.T().Context())
	}
	if s.cockroachContainer != nil {
		_ = s.cockroachContainer.Terminate(s.T().Context())
	}
}

func (s *SSOServiceTestSuite) SetupTest() {
	ctx := context.Background()

	_, err := s.db.Exec(ctx, "DELETE FROM tokens")
	require.NoError(s.T(), err, "failed to clean tokens table")

	keys, err := s.rdb.Keys(ctx, "jwks:*").Result()
	if err == nil && len(keys) > 0 {
		_ = s.rdb.Del(ctx, keys...).Err()
	}
}

func (s *SSOServiceTestSuite) TearDownTest() {
	ctx := context.Background()

	_, err := s.db.Exec(ctx, "DELETE FROM tokens")
	require.NoError(s.T(), err, "failed to clean tokens table")

	keys, err := s.rdb.Keys(ctx, "jwks:*").Result()
	if err == nil && len(keys) > 0 {
		_ = s.rdb.Del(ctx, keys...).Err()
	}
}

func (s *SSOServiceTestSuite) TestGetToken_PasswordGrant() {
	s.Run("should get token with password grant successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		email := "test@example.com"
		password := "password123"
		clientID := "client-123"
		userID := uuid.New().String()

		user := usermodel.User{
			ID:       userID,
			Email:    email,
			Password: password,
			Role:     usermodel.RoleUser,
		}

		s.userRepo.EXPECT().
			GetByEmail(gomock.Any(), email).
			Return(user, nil)

		keyRepo, err := redis.New(ctx, s.rdb)
		require.NoError(s.T(), err)
		_, err = keyRepo.Rotate(ctx)
		require.NoError(s.T(), err)

		req := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_PasswordGrant{
				PasswordGrant: &pb.PasswordGrant{
					Username: email,
					Password: password,
				},
			},
			ClientId: clientID,
		}

		resp, err := s.ssoClient.GetToken(ctx, req)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), resp.Token)
		assert.NotEmpty(s.T(), resp.Token.AccessToken)
		assert.NotEmpty(s.T(), resp.Token.RefreshToken)
		assert.Equal(s.T(), pb.TokenType_TOKEN_TYPE_BEARER, resp.Token.TokenType)
	})

	s.Run("should return error when user not found", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		email := "notfound@example.com"

		s.userRepo.EXPECT().
			GetByEmail(gomock.Any(), email).
			Return(usermodel.User{}, errors.New("user not found"))

		req := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_PasswordGrant{
				PasswordGrant: &pb.PasswordGrant{
					Username: email,
					Password: "password123",
				},
			},
			ClientId: "client-123",
		}

		resp, err := s.ssoClient.GetToken(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should return error when password is invalid", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		email := "invalidpass@example.com"
		password := "correct-password"
		wrongPassword := "wrong-password"
		userID := uuid.New().String()

		user := usermodel.User{
			ID:       userID,
			Email:    email,
			Password: password,
			Role:     usermodel.RoleUser,
		}

		s.userRepo.EXPECT().
			GetByEmail(gomock.Any(), email).
			Return(user, nil)

		req := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_PasswordGrant{
				PasswordGrant: &pb.PasswordGrant{
					Username: email,
					Password: wrongPassword,
				},
			},
			ClientId: "client-123",
		}

		resp, err := s.ssoClient.GetToken(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should return error when no key exists", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		keys, err := s.rdb.Keys(ctx, "jwks:*").Result()
		if err == nil && len(keys) > 0 {
			_ = s.rdb.Del(ctx, keys...).Err()
		}

		email := "nokey@example.com"
		password := "password123"
		userID := uuid.New().String()

		user := usermodel.User{
			ID:       userID,
			Email:    email,
			Password: password,
			Role:     usermodel.RoleUser,
		}

		s.userRepo.EXPECT().
			GetByEmail(gomock.Any(), email).
			Return(user, nil)

		req := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_PasswordGrant{
				PasswordGrant: &pb.PasswordGrant{
					Username: email,
					Password: password,
				},
			},
			ClientId: "client-123",
		}

		resp, err := s.ssoClient.GetToken(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})
}

func (s *SSOServiceTestSuite) TestGetToken_RefreshTokenGrant() {
	s.Run("should get token with refresh token grant successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		email := "refresh@example.com"
		password := "password123"
		clientID := "client-123"
		userID := uuid.New().String()

		user := usermodel.User{
			ID:       userID,
			Email:    email,
			Password: password,
			Role:     usermodel.RoleUser,
		}

		s.userRepo.EXPECT().
			GetByEmail(gomock.Any(), email).
			Return(user, nil)

		keyRepo, err := redis.New(ctx, s.rdb)
		require.NoError(s.T(), err)
		_, err = keyRepo.Rotate(ctx)
		require.NoError(s.T(), err)

		passwordReq := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_PasswordGrant{
				PasswordGrant: &pb.PasswordGrant{
					Username: email,
					Password: password,
				},
			},
			ClientId: clientID,
		}

		passwordResp, err := s.ssoClient.GetToken(ctx, passwordReq)
		require.NoError(s.T(), err)
		require.NotNil(s.T(), passwordResp.Token)
		require.NotEmpty(s.T(), passwordResp.Token.RefreshToken)

		sum := sha256.Sum256([]byte(passwordResp.Token.RefreshToken))
		tokenHash := hex.EncodeToString(sum[:])

		refreshReq := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_RefreshTokenGrant{
				RefreshTokenGrant: &pb.RefreshTokenGrant{
					RefreshToken: tokenHash,
				},
			},
			ClientId: clientID,
		}

		refreshResp, err := s.ssoClient.GetToken(ctx, refreshReq)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), refreshResp.Token)
		assert.NotEmpty(s.T(), refreshResp.Token.AccessToken)
		assert.NotEmpty(s.T(), refreshResp.Token.RefreshToken)
		assert.Equal(s.T(), pb.TokenType_TOKEN_TYPE_BEARER, refreshResp.Token.TokenType)
		assert.NotEqual(s.T(), passwordResp.Token.RefreshToken, refreshResp.Token.RefreshToken)
	})

	s.Run("should return error when refresh token is invalid", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_RefreshTokenGrant{
				RefreshTokenGrant: &pb.RefreshTokenGrant{
					RefreshToken: "invalid-refresh-token",
				},
			},
			ClientId: "client-123",
		}

		resp, err := s.ssoClient.GetToken(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should return error when refresh token is reused", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		email := "reuse@example.com"
		password := "password123"
		clientID := "client-123"
		userID := uuid.New().String()

		user := usermodel.User{
			ID:       userID,
			Email:    email,
			Password: password,
			Role:     usermodel.RoleUser,
		}

		s.userRepo.EXPECT().
			GetByEmail(gomock.Any(), email).
			Return(user, nil)

		keyRepo, err := redis.New(ctx, s.rdb)
		require.NoError(s.T(), err)
		_, err = keyRepo.Rotate(ctx)
		require.NoError(s.T(), err)

		passwordReq := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_PasswordGrant{
				PasswordGrant: &pb.PasswordGrant{
					Username: email,
					Password: password,
				},
			},
			ClientId: clientID,
		}

		passwordResp, err := s.ssoClient.GetToken(ctx, passwordReq)
		require.NoError(s.T(), err)
		refreshToken := passwordResp.Token.RefreshToken

		sum := sha256.Sum256([]byte(refreshToken))
		tokenHash := hex.EncodeToString(sum[:])

		refreshReq := &pb.GetToken_Request{
			Grant: &pb.GetToken_Request_RefreshTokenGrant{
				RefreshTokenGrant: &pb.RefreshTokenGrant{
					RefreshToken: tokenHash,
				},
			},
			ClientId: clientID,
		}

		_, err = s.ssoClient.GetToken(ctx, refreshReq)
		require.NoError(s.T(), err)

		_, err = s.ssoClient.GetToken(ctx, refreshReq)
		assert.Error(s.T(), err)
	})
}

func (s *SSOServiceTestSuite) TestGetJWKS() {
	s.Run("should get JWKS successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		keyRepo, err := redis.New(ctx, s.rdb)
		require.NoError(s.T(), err)

		key1, err := keyRepo.Rotate(ctx)
		require.NoError(s.T(), err)

		time.Sleep(time.Second)
		key2, err := keyRepo.Rotate(ctx)
		require.NoError(s.T(), err)

		req := &pb.GetJWKS_Request{}

		resp, err := s.ssoClient.GetJWKS(ctx, req)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), resp.Keys)
		assert.Len(s.T(), resp.Keys, 2)

		keyIDs := make(map[string]bool)
		for _, key := range resp.Keys {
			keyIDs[key.Kid] = true
			assert.NotEmpty(s.T(), key.Kid)
			assert.NotEmpty(s.T(), key.Public)
			assert.Equal(s.T(), "RS256", key.Algorithm)
			assert.True(s.T(), key.Active)
		}

		assert.True(s.T(), keyIDs[key1.ID])
		assert.True(s.T(), keyIDs[key2.ID])
	})

	s.Run("should return empty JWKS when no keys exist", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		keys, err := s.rdb.Keys(ctx, "jwks:*").Result()
		if err == nil && len(keys) > 0 {
			_ = s.rdb.Del(ctx, keys...).Err()
		}

		req := &pb.GetJWKS_Request{}

		resp, err := s.ssoClient.GetJWKS(ctx, req)
		require.NoError(s.T(), err)
		require.NotNil(s.T(), resp)
		assert.True(s.T(), len(resp.Keys) == 0)
	})
}

func (s *SSOServiceTestSuite) TestGetToken_InvalidGrant() {
	s.Run("should return error when grant type is invalid", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.GetToken_Request{
			ClientId: "client-123",
		}

		resp, err := s.ssoClient.GetToken(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestSSOServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SSOServiceTestSuite))
}

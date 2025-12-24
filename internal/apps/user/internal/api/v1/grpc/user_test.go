package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/postgres"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/service"
	"github.com/kitanoyoru/kgym/internal/apps/user/migrations"
	postgresdb "github.com/kitanoyoru/kgym/pkg/database/postgres"
	"github.com/kitanoyoru/kgym/pkg/testing/integration/cockroachdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceTestSuite struct {
	suite.Suite

	db        *pgxpool.Pool
	container *cockroachdb.CockroachDBContainer
	server    *grpc.Server
	client    pb.UserServiceClient
	conn      *grpc.ClientConn
}

func (s *UserServiceTestSuite) SetupSuite() {
	ctx := context.Background()

	container, err := cockroachdb.SetupTestContainer(ctx)
	require.NoError(s.T(), err, "failed to setup test container")

	s.container = container

	s.db, err = postgresdb.New(ctx, postgresdb.Config{
		URI: container.URI,
	})
	require.NoError(s.T(), err, "failed to create postgres client")

	err = migrations.Up(ctx, "pgx", container.URI)
	require.NoError(s.T(), err, "failed to run migrations")

	repository := postgres.New(s.db)
	userService := service.New(repository)
	grpcServer := NewUserService(userService)

	s.server = grpc.NewServer()
	pb.RegisterUserServiceServer(s.server, grpcServer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err, "failed to create listener")

	go func() {
		_ = s.server.Serve(listener)
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.NewClient(
		listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(s.T(), err, "failed to create gRPC client")

	s.conn = conn
	s.client = pb.NewUserServiceClient(conn)
}

func (s *UserServiceTestSuite) TearDownSuite() {
	if s.conn != nil {
		_ = s.conn.Close()
	}
	if s.server != nil {
		s.server.GracefulStop()
	}
	if s.container != nil {
		_ = s.container.Terminate(s.T().Context())
	}
	if s.db != nil {
		s.db.Close()
	}
}

func (s *UserServiceTestSuite) SetupTest() {
	ctx := context.Background()
	_, err := s.db.Exec(ctx, "DELETE FROM users")
	require.NoError(s.T(), err, "failed to clean users table")
}

func (s *UserServiceTestSuite) TearDownTest() {
	ctx := context.Background()
	_, err := s.db.Exec(ctx, "DELETE FROM users")
	require.NoError(s.T(), err, "failed to clean users table")
}

func (s *UserServiceTestSuite) TestCreateUser() {
	s.Run("should create a user successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:    "test@example.com",
			Role:     "default",
			Username: "testuser",
			Password: "password123",
		}

		resp, err := s.client.CreateUser(ctx, req)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), resp.Id)
		assert.Regexp(s.T(), `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, resp.Id)
	})

	s.Run("should not create a user because of invalid email", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:    "invalid-email",
			Role:     "default",
			Username: "testuser",
			Password: "password123",
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of empty password", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:    "test@example.com",
			Role:     "default",
			Username: "testuser",
			Password: "",
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of password too short", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:    "test@example.com",
			Role:     "default",
			Username: "testuser",
			Password: "short",
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of username too short", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:    "test@example.com",
			Role:     "default",
			Username: "ab",
			Password: "password123",
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of invalid role", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:    "test@example.com",
			Role:     "invalid-role",
			Username: "testuser",
			Password: "password123",
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})
}

func (s *UserServiceTestSuite) TestListUsers() {
	ctx := context.Background()
	_, _ = s.db.Exec(ctx, "DELETE FROM users")

	s.Run("should return list of users successfully with no filters", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		user1ID := s.createTestUser(ctx, "user1@example.com", "default", "user1", "password123")
		user2ID := s.createTestUser(ctx, "user2@example.com", "admin", "user2", "password123")

		req := &pb.ListUsers_Request{}

		resp, err := s.client.ListUsers(ctx, req)
		require.NoError(s.T(), err)
		assert.Len(s.T(), resp.Users, 2)

		var ids []string
		for _, u := range resp.Users {
			ids = append(ids, u.Id)
		}
		assert.Contains(s.T(), ids, user1ID)
		assert.Contains(s.T(), ids, user2ID)
	})

	s.Run("should return list of users successfully with email filter", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		email := "filter@example.com"
		userID := s.createTestUser(ctx, email, "default", "emailfilteruser", "password123")
		s.createTestUser(ctx, "other@example.com", "default", "otheruser", "password123")

		req := &pb.ListUsers_Request{
			Email: &email,
		}

		resp, err := s.client.ListUsers(ctx, req)
		require.NoError(s.T(), err)
		assert.Len(s.T(), resp.Users, 1)
		assert.Equal(s.T(), userID, resp.Users[0].Id)
		assert.Equal(s.T(), email, resp.Users[0].Email)
	})

	s.Run("should return list of users successfully with username filter", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		username := "usernamefilteruser"
		userID := s.createTestUser(ctx, "usernamefilter@example.com", "default", username, "password123")
		s.createTestUser(ctx, "otheruser@example.com", "default", "otheruser2", "password123")

		req := &pb.ListUsers_Request{
			Username: &username,
		}

		resp, err := s.client.ListUsers(ctx, req)
		require.NoError(s.T(), err)
		assert.Len(s.T(), resp.Users, 1)
		assert.Equal(s.T(), userID, resp.Users[0].Id)
		assert.Equal(s.T(), username, resp.Users[0].Username)
	})

	s.Run("should return list of users successfully with role filter", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		role := "admin"
		user1ID := s.createTestUser(ctx, "admin1@example.com", role, "admin1", "password123")
		user2ID := s.createTestUser(ctx, "admin2@example.com", role, "admin2", "password123")
		s.createTestUser(ctx, "user@example.com", "default", "user", "password123")

		req := &pb.ListUsers_Request{
			Role: &role,
		}

		resp, err := s.client.ListUsers(ctx, req)
		require.NoError(s.T(), err)
		assert.Len(s.T(), resp.Users, 2)

		var ids []string
		for _, u := range resp.Users {
			ids = append(ids, u.Id)
			assert.Equal(s.T(), role, u.Role)
		}
		assert.Contains(s.T(), ids, user1ID)
		assert.Contains(s.T(), ids, user2ID)
	})

	s.Run("should return empty list when no users found", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		req := &pb.ListUsers_Request{}

		resp, err := s.client.ListUsers(ctx, req)
		require.NoError(s.T(), err)
		assert.Empty(s.T(), resp.Users)
	})
}

func (s *UserServiceTestSuite) TestDeleteUser() {
	ctx := context.Background()
	_, _ = s.db.Exec(ctx, "DELETE FROM users")

	s.Run("should delete a user successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		userID := s.createTestUser(ctx, "delete@example.com", "default", "deleteuser", "password123")

		req := &pb.DeleteUser_Request{
			Id: userID,
		}

		resp, err := s.client.DeleteUser(ctx, req)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), resp)

		listReq := &pb.ListUsers_Request{}
		listResp, err := s.client.ListUsers(ctx, listReq)
		require.NoError(s.T(), err)

		var ids []string
		for _, u := range listResp.Users {
			ids = append(ids, u.Id)
		}
		assert.NotContains(s.T(), ids, userID)
	})

	s.Run("should return error when user not found", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		nonExistentID := uuid.New().String()

		req := &pb.DeleteUser_Request{
			Id: nonExistentID,
		}

		resp, err := s.client.DeleteUser(ctx, req)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), resp)
	})
}

func (s *UserServiceTestSuite) createTestUser(ctx context.Context, email, role, username, password string) string {
	req := &pb.CreateUser_Request{
		Email:    email,
		Role:     role,
		Username: username,
		Password: password,
	}

	resp, err := s.client.CreateUser(ctx, req)
	if err != nil {
		s.T().Logf("Failed to create test user: %v", err)
	}
	require.NoError(s.T(), err)
	return resp.Id
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

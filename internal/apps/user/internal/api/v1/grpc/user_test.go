package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/repository/user/postgres"
	userservice "github.com/kitanoyoru/kgym/internal/apps/user/internal/service/user"
	"github.com/kitanoyoru/kgym/internal/apps/user/migrations"
	postgresdb "github.com/kitanoyoru/kgym/pkg/database/postgres"
	"github.com/kitanoyoru/kgym/pkg/testing/integration/cockroachdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	userService := userservice.New(repository)
	grpcServer, err := NewUserService(userService)
	require.NoError(s.T(), err, "failed to create gRPC server")

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

		birthDate := carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()
		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role_USER,
			Username:  "testuser",
			Password:  "password123",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: timestamppb.New(birthDate),
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
			Email:     "invalid-email",
			Role:      pb.Role_USER,
			Username:  "testuser",
			Password:  "password123",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of empty password", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role_USER,
			Username:  "testuser",
			Password:  "",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of password too short", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role_USER,
			Username:  "testuser",
			Password:  "short",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of username too short", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role_USER,
			Username:  "ab",
			Password:  "password123",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of invalid role", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role(999),
			Username:  "testuser",
			Password:  "password123",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of invalid avatar URL", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role_USER,
			Username:  "testuser",
			Password:  "password123",
			AvatarUrl: "not-a-valid-url",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of invalid mobile", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role_USER,
			Username:  "testuser",
			Password:  "password123",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "123",
			FirstName: "John",
			LastName:  "Doe",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of empty first name", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role_USER,
			Username:  "testuser",
			Password:  "password123",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "",
			LastName:  "Doe",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should not create a user because of empty last name", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CreateUser_Request{
			Email:     "test@example.com",
			Role:      pb.Role_USER,
			Username:  "testuser",
			Password:  "password123",
			AvatarUrl: "https://example.com/avatar.jpg",
			Mobile:    "+1234567890",
			FirstName: "John",
			LastName:  "",
			BirthDate: timestamppb.New(carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime()),
		}

		resp, err := s.client.CreateUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})
}

func (s *UserServiceTestSuite) TestGetUser() {
	ctx := context.Background()
	_, _ = s.db.Exec(ctx, "DELETE FROM users")

	s.Run("should get user by id successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		email := "getbyid@example.com"
		userID := s.createTestUser(ctx, email, pb.Role_USER, "getbyiduser", "password123")

		req := &pb.GetUser_Request{
			Id: &userID,
		}

		resp, err := s.client.GetUser(ctx, req)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), resp.User)
		assert.Equal(s.T(), userID, resp.User.Id)
		assert.Equal(s.T(), email, resp.User.Email)
	})

	s.Run("should get user by email successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		email := "getbyemail@example.com"
		userID := s.createTestUser(ctx, email, pb.Role_USER, "getbyemailuser", "password123")

		req := &pb.GetUser_Request{
			Email: &email,
		}

		resp, err := s.client.GetUser(ctx, req)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), resp.User)
		assert.Equal(s.T(), userID, resp.User.Id)
		assert.Equal(s.T(), email, resp.User.Email)
	})

	s.Run("should return error when neither id nor email provided", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.GetUser_Request{}

		resp, err := s.client.GetUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should return error when user not found by id", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		nonExistentID := uuid.New().String()
		req := &pb.GetUser_Request{
			Id: &nonExistentID,
		}

		resp, err := s.client.GetUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})

	s.Run("should return error when user not found by email", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		nonExistentEmail := "nonexistent@example.com"
		req := &pb.GetUser_Request{
			Email: &nonExistentEmail,
		}

		resp, err := s.client.GetUser(ctx, req)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), resp)
	})
}

func (s *UserServiceTestSuite) TestDeleteUser() {
	ctx := context.Background()
	_, _ = s.db.Exec(ctx, "DELETE FROM users")

	s.Run("should delete a user successfully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.db.Exec(ctx, "DELETE FROM users")

		userID := s.createTestUser(ctx, "delete@example.com", pb.Role_USER, "deleteuser", "password123")

		req := &pb.DeleteUser_Request{
			Id: userID,
		}

		resp, err := s.client.DeleteUser(ctx, req)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), resp)

		getReq := &pb.GetUser_Request{
			Id: &userID,
		}
		_, err = s.client.GetUser(ctx, getReq)
		assert.Error(s.T(), err)
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

func (s *UserServiceTestSuite) createTestUser(ctx context.Context, email string, role pb.Role, username, password string) string {
	return s.createTestUserWithFields(ctx, email, role, username, password, "https://example.com/avatar.jpg", "+1234567890", "John", "Doe", carbon.CreateFromDateTime(1990, 1, 1, 0, 0, 0).SetTimezone(carbon.UTC).StdTime())
}

func (s *UserServiceTestSuite) createTestUserWithFields(ctx context.Context, email string, role pb.Role, username, password, avatarURL, mobile, firstName, lastName string, birthDate time.Time) string {
	req := &pb.CreateUser_Request{
		Email:     email,
		Role:      role,
		Username:  username,
		Password:  password,
		AvatarUrl: avatarURL,
		Mobile:    mobile,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: timestamppb.New(birthDate),
	}

	resp, err := s.client.CreateUser(ctx, req)
	if err != nil {
		s.T().Logf("Failed to create test user: %v", err)
	}
	require.NoError(s.T(), err)
	return resp.Id
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

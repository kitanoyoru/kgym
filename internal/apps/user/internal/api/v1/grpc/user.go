package grpc

import (
	"context"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/service"
	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics"
	apimetrics "github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	service *service.Service
}

func NewUserService(service *service.Service) *UserServiceServer {
	return &UserServiceServer{
		service: service,
	}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUser_Request) (*pb.CreateUser_Response, error) {
	metrics.GlobalRegistry.GetMetric(apimetrics.GRPCMethodCreateUser).Counter.WithLabelValues().Inc()

	id, err := s.service.Create(ctx, service.CreateUserRequest{
		Email:    req.Email,
		Role:     req.Role,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUser_Response{
		Id: id,
	}, nil
}

func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsers_Request) (*pb.ListUsers_Response, error) {
	metrics.GlobalRegistry.GetMetric(apimetrics.GRPCMethodListUsers).Counter.WithLabelValues().Inc()

	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUser_Request) (*pb.DeleteUser_Response, error) {
	metrics.GlobalRegistry.GetMetric(apimetrics.GRPCMethodDeleteUser).Counter.WithLabelValues().Inc()

	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

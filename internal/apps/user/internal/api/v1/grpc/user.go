package grpc

import (
	"context"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/api/v1/grpc/serializer"
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

	id, err := s.service.Create(ctx, serializer.PbCreateRequestToServiceRequest(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUser_Response{
		Id: id,
	}, nil
}

func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsers_Request) (*pb.ListUsers_Response, error) {
	metrics.GlobalRegistry.GetMetric(apimetrics.GRPCMethodListUsers).Counter.WithLabelValues().Inc()

	options := serializer.PbListRequestToServiceOptions(req)

	users, err := s.service.List(ctx, options...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = serializer.EntityToPbUser(&user)
	}

	return &pb.ListUsers_Response{
		Users: pbUsers,
	}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUser_Request) (*pb.DeleteUser_Response, error) {
	metrics.GlobalRegistry.GetMetric(apimetrics.GRPCMethodDeleteUser).Counter.WithLabelValues().Inc()

	err := s.service.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUser_Response{}, nil
}

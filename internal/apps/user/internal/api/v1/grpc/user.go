package grpc

import (
	"context"
	"fmt"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/api/v1/grpc/serializer"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/service"
	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics"
	"github.com/kitanoyoru/kgym/pkg/metrics/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	GRPCServiceMetricsPrefix = "kgym.user.api.grpc"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	service *service.Service
}

func NewUserService(service *service.Service) (*UserServiceServer, error) {
	methods := []string{
		"CreateUser",
		"ListUsers",
		"DeleteUser",
	}

	for _, method := range methods {
		if err := metrics.GlobalRegistry.RegisterMetric(prometheus.MetricConfig{
			Name: fmt.Sprintf("%s.%s", GRPCServiceMetricsPrefix, method),
			Type: prometheus.Counter,
		}); err != nil {
			return nil, err
		}
	}

	return &UserServiceServer{
		service: service,
	}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUser_Request) (*pb.CreateUser_Response, error) {
	id, err := s.service.Create(ctx, serializer.PbCreateRequestToServiceRequest(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUser_Response{
		Id: id,
	}, nil
}

func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsers_Request) (*pb.ListUsers_Response, error) {
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
	err := s.service.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUser_Response{}, nil
}

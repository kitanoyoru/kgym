package grpc

import (
	"context"
	"fmt"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/api/v1/grpc/serializer"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	userservice "github.com/kitanoyoru/kgym/internal/apps/user/internal/service/user"
	"github.com/kitanoyoru/kgym/pkg/metrics/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	GRPCServiceMetricsPrefix = "kgym.user.api.grpc"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	svc userservice.IService
}

func NewUserService(svc userservice.IService) (*UserServiceServer, error) {
	methods := []string{
		"CreateUser",
		"GetUser",
		"DeleteUser",
	}

	for _, method := range methods {
		if err := prometheus.GlobalRegistry.RegisterMetric(prometheus.MetricConfig{
			Name: fmt.Sprintf("%s.%s", GRPCServiceMetricsPrefix, method),
			Type: prometheus.Counter,
		}); err != nil {
			return nil, err
		}
	}

	return &UserServiceServer{
		svc: svc,
	}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUser_Request) (*pb.CreateUser_Response, error) {
	svcReq, err := serializer.PbCreateRequestToServiceRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	svcResp, err := s.svc.Create(ctx, svcReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUser_Response{
		Id: svcResp.ID,
	}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUser_Request) (*pb.GetUser_Response, error) {
	var userEntity userentity.User
	var err error

	id := req.GetId()
	email := req.GetEmail()

	if id != "" {
		userEntity, err = s.svc.GetByID(ctx, id)
	} else if email != "" {
		userEntity, err = s.svc.GetByEmail(ctx, email)
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "either id or email must be provided")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	pbUser, err := serializer.EntityToPbUser(userEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert entity to protobuf: %v", err)
	}

	return &pb.GetUser_Response{
		User: &pbUser,
	}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUser_Request) (*pb.DeleteUser_Response, error) {
	if err := s.svc.DeleteByID(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUser_Response{}, nil
}

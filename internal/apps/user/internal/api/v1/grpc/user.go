package grpc

import (
	"context"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/api/v1/grpc/serializer"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	userservice "github.com/kitanoyoru/kgym/internal/apps/user/internal/service/user"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	GRPCServicePrefix = "kgym.user.api.grpc"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer

	tracer      trace.Tracer
	userService userservice.IService
}

func NewUserService(userService userservice.IService) (*UserServiceServer, error) {
	tracer := otel.Tracer(GRPCServicePrefix)

	return &UserServiceServer{
		tracer:      tracer,
		userService: userService,
	}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUser_Request) (*pb.CreateUser_Response, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUser")
	defer span.End()

	svcReq, err := serializer.PbCreateRequestToServiceRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	svcResp, err := s.userService.Create(ctx, svcReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUser_Response{
		Id: svcResp.ID,
	}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUser_Request) (*pb.GetUser_Response, error) {
	ctx, span := s.tracer.Start(ctx, "GetUser")
	defer span.End()

	var userEntity userentity.User
	var err error

	id := req.GetId()
	email := req.GetEmail()

	if id != "" {
		userEntity, err = s.userService.GetByID(ctx, id)
	} else if email != "" {
		userEntity, err = s.userService.GetByEmail(ctx, email)
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
	ctx, span := s.tracer.Start(ctx, "DeleteUser")
	defer span.End()

	if err := s.userService.DeleteByID(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUser_Response{}, nil
}

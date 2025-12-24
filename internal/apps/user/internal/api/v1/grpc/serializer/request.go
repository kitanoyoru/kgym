package serializer

import (
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/service"
)

func PbCreateRequestToServiceRequest(pbCreateRequest *pb.CreateUser_Request) service.CreateUserRequest {
	return service.CreateUserRequest{
		Email:    pbCreateRequest.Email,
		Role:     pbCreateRequest.Role,
		Username: pbCreateRequest.Username,
		Password: pbCreateRequest.Password,
	}
}

func PbListRequestToServiceOptions(pbListRequest *pb.ListUsers_Request) []service.Option {
	var options []service.Option
	if pbListRequest.Email != nil {
		options = append(options, service.WithEmail(*pbListRequest.Email))
	}
	if pbListRequest.Role != nil {
		options = append(options, service.WithRole(userentity.Role(*pbListRequest.Role)))
	}
	if pbListRequest.Username != nil {
		options = append(options, service.WithUsername(*pbListRequest.Username))
	}
	return options
}

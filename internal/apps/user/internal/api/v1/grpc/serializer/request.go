package serializer

import (
	"github.com/dromara/carbon/v2"
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	"github.com/kitanoyoru/kgym/internal/apps/user/internal/service"
)

func PbCreateRequestToServiceRequest(pbCreateRequest *pb.CreateUser_Request) service.CreateUserRequest {
	return service.CreateUserRequest{
		Email:     pbCreateRequest.Email,
		Role:      roleProtoToString(pbCreateRequest.Role),
		Username:  pbCreateRequest.Username,
		Password:  pbCreateRequest.Password,
		AvatarURL: pbCreateRequest.AvatarUrl,
		Mobile:    pbCreateRequest.Mobile,
		FirstName: pbCreateRequest.FirstName,
		LastName:  pbCreateRequest.LastName,
		BirthDate: carbon.CreateFromStdTime(pbCreateRequest.BirthDate.AsTime()).SetTimezone(carbon.UTC).StdTime(),
	}
}

func PbListRequestToServiceOptions(pbListRequest *pb.ListUsers_Request) []service.Option {
	var options []service.Option
	if pbListRequest.Email != nil {
		options = append(options, service.WithEmail(*pbListRequest.Email))
	}
	if pbListRequest.Role != nil {
		options = append(options, service.WithRole(userentity.Role(roleProtoToString(*pbListRequest.Role))))
	}
	if pbListRequest.Username != nil {
		options = append(options, service.WithUsername(*pbListRequest.Username))
	}
	return options
}

func roleProtoToString(role pb.Role) string {
	switch role {
	case pb.Role_ADMIN:
		return "admin"
	case pb.Role_USER:
		return "default"
	default:
		return "default"
	}
}

func roleStringToProto(role string) pb.Role {
	switch role {
	case "admin":
		return pb.Role_ADMIN
	case "default":
		return pb.Role_USER
	default:
		return pb.Role_USER
	}
}

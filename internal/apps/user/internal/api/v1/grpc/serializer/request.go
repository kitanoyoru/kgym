package serializer

import (
	"github.com/dromara/carbon/v2"
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	userservice "github.com/kitanoyoru/kgym/internal/apps/user/internal/service"
)

func PbCreateRequestToServiceRequest(pbCreateRequest *pb.CreateUser_Request) (userservice.CreateRequest, error) {
	role, err := userentity.RoleFromString(pbCreateRequest.Role.String())
	if err != nil {
		return userservice.CreateRequest{}, err
	}

	return userservice.CreateRequest{
		Email:     pbCreateRequest.Email,
		Role:      role,
		Username:  pbCreateRequest.Username,
		Password:  pbCreateRequest.Password,
		AvatarURL: pbCreateRequest.AvatarUrl,
		Mobile:    pbCreateRequest.Mobile,
		FirstName: pbCreateRequest.FirstName,
		LastName:  pbCreateRequest.LastName,
		BirthDate: carbon.CreateFromStdTime(pbCreateRequest.BirthDate.AsTime()).SetTimezone(carbon.UTC).StdTime(),
	}, nil
}

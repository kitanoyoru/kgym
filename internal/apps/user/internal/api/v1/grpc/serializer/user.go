package serializer

import (
	"github.com/dromara/carbon/v2"
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PbUserToEntity(pbUser *pb.User) (*userentity.User, error) {
	role, err := pbRoleToEntity(pbUser.Role)
	if err != nil {
		return nil, err
	}

	return &userentity.User{
		ID:        pbUser.Id,
		Email:     pbUser.Email,
		Role:      role,
		Username:  pbUser.Username,
		Password:  pbUser.Password,
		AvatarURL: pbUser.AvatarUrl,
		Mobile:    pbUser.Mobile,
		FirstName: pbUser.FirstName,
		LastName:  pbUser.LastName,
		BirthDate: carbon.CreateFromStdTime(pbUser.BirthDate.AsTime()).SetTimezone(carbon.UTC).StdTime(),
	}, nil
}

func EntityToPbUser(entityUser userentity.User) (pb.User, error) {
	role, err := entityRoleToProto(entityUser.Role)
	if err != nil {
		return pb.User{}, err
	}

	return pb.User{
		Id:        entityUser.ID,
		Email:     entityUser.Email,
		Role:      role,
		Username:  entityUser.Username,
		Password:  entityUser.Password,
		AvatarUrl: entityUser.AvatarURL,
		Mobile:    entityUser.Mobile,
		FirstName: entityUser.FirstName,
		LastName:  entityUser.LastName,
		BirthDate: timestamppb.New(entityUser.BirthDate),
	}, nil
}

func entityRoleToProto(role userentity.Role) (pb.Role, error) {
	switch role {
	case userentity.RoleAdmin:
		return pb.Role_ADMIN, nil
	case userentity.RoleUser:
		return pb.Role_USER, nil
	default:
		return pb.Role_USER, errors.New("invalid role")
	}
}

func pbRoleToEntity(role pb.Role) (userentity.Role, error) {
	switch role {
	case pb.Role_ADMIN:
		return userentity.RoleAdmin, nil
	case pb.Role_USER:
		return userentity.RoleUser, nil
	default:
		return userentity.RoleUser, errors.New("invalid role")
	}
}

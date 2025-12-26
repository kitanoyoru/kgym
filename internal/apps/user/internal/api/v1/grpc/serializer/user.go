package serializer

import (
	"github.com/dromara/carbon/v2"
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PbUserToEntity(pbUser *pb.User) *userentity.User {
	return &userentity.User{
		ID:        pbUser.Id,
		Email:     pbUser.Email,
		Role:      userentity.Role(roleProtoToString(pbUser.Role)),
		Username:  pbUser.Username,
		Password:  pbUser.Password,
		AvatarURL: pbUser.AvatarUrl,
		Mobile:    pbUser.Mobile,
		FirstName: pbUser.FirstName,
		LastName:  pbUser.LastName,
		BirthDate: carbon.CreateFromStdTime(pbUser.BirthDate.AsTime()).SetTimezone(carbon.UTC).StdTime(),
	}
}

func EntityToPbUser(entityUser *userentity.User) *pb.User {
	return &pb.User{
		Id:        entityUser.ID,
		Email:     entityUser.Email,
		Role:      roleStringToProto(string(entityUser.Role)),
		Username:  entityUser.Username,
		Password:  entityUser.Password,
		AvatarUrl: entityUser.AvatarURL,
		Mobile:    entityUser.Mobile,
		FirstName: entityUser.FirstName,
		LastName:  entityUser.LastName,
		BirthDate: timestamppb.New(entityUser.BirthDate),
	}
}

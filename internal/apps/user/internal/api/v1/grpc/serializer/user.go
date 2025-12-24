package serializer

import (
	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
)

func PbUserToEntity(pbUser *pb.User) *userentity.User {
	return &userentity.User{
		ID:       pbUser.Id,
		Email:    pbUser.Email,
		Role:     userentity.Role(pbUser.Role),
		Username: pbUser.Username,
		Password: pbUser.Password,
	}
}

func EntityToPbUser(entityUser *userentity.User) *pb.User {
	return &pb.User{
		Id:       entityUser.ID,
		Email:    entityUser.Email,
		Role:     string(entityUser.Role),
		Username: entityUser.Username,
		Password: entityUser.Password,
	}
}

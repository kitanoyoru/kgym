package user

import (
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
	"github.com/pkg/errors"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func RoleFromEntity(entity userentity.Role) (Role, error) {
	switch entity {
	case userentity.RoleAdmin:
		return RoleAdmin, nil
	case userentity.RoleUser:
		return RoleUser, nil
	default:
		return "", errors.New("invalid role")
	}
}

func (r Role) ToEntity() (userentity.Role, error) {
	switch r {
	case RoleAdmin:
		return userentity.RoleAdmin, nil
	case RoleUser:
		return userentity.RoleUser, nil
	default:
		return "", errors.New("invalid role")
	}
}

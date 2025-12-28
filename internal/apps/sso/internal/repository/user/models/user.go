package models

import (
	"time"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/user/v1"
	"github.com/pkg/errors"
)

var (
	ErrRoleNotFound = errors.New("role not found")
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func RoleFromProto(role pb.Role) (Role, error) {
	switch role {
	case pb.Role_ADMIN:
		return RoleAdmin, nil
	case pb.Role_USER:
		return RoleUser, nil
	default:
		return "", ErrRoleNotFound
	}
}

type User struct {
	ID        string
	Email     string
	Role      Role
	Username  string
	Password  string
	AvatarURL string
	Mobile    string
	FirstName string
	LastName  string
	BirthDate time.Time
	CreatedAt time.Time
}

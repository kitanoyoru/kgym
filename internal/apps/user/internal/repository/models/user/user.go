package user

import (
	"time"

	"github.com/dromara/carbon/v2"
	userentity "github.com/kitanoyoru/kgym/internal/apps/user/internal/entity/user"
)

const (
	Table = "users"
)

var Columns = []string{
	"id",
	"email",
	"role",
	"username",
	"password",
	"avatar_url",
	"mobile",
	"first_name",
	"last_name",
	"birth_date",
	"created_at",
	"updated_at",
	"deleted_at",
}

func UserFromEntity(entity userentity.User) (User, error) {
	role, err := RoleFromEntity(entity.Role)
	if err != nil {
		return User{}, err
	}

	now := carbon.Now().StdTime()

	return User{
		ID:        entity.ID,
		Email:     entity.Email,
		Role:      role,
		Username:  entity.Username,
		Password:  entity.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

type User struct {
	ID        string     `db:"id"`
	Email     string     `db:"email"`
	Role      Role       `db:"role"`
	Username  string     `db:"username"`
	Password  string     `db:"password"`
	AvatarURL string     `db:"avatar_url"`
	Mobile    string     `db:"mobile"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
	BirthDate time.Time  `db:"birth_date"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (u User) Values() []any {
	return []any{
		u.ID,
		u.Email,
		u.Role,
		u.Username,
		u.Password,
		u.AvatarURL,
		u.Mobile,
		u.FirstName,
		u.LastName,
		u.BirthDate,
		u.CreatedAt,
		u.UpdatedAt,
		u.DeletedAt,
	}
}

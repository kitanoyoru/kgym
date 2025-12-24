package models

import (
	"time"
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
	"created_at",
	"updated_at",
	"deleted_at",
}

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleDefault Role = "default"
)

type User struct {
	ID        string     `db:"id"`
	Email     string     `db:"email"`
	Role      Role       `db:"role"`
	Username  string     `db:"username"`
	Password  string     `db:"password"`
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
		u.CreatedAt,
		u.UpdatedAt,
		u.DeletedAt,
	}
}
